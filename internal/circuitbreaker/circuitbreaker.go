package circuitbreaker

import (
	"context"
	"errors"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/sony/gobreaker"
)

type CircuitBreakerConfig struct {
	FailureThreshold int           // Number of failures before opening circuit (default: 5)
	SuccessThreshold int           // Number of successes needed to close circuit (default: 3)
	Timeout          time.Duration // How long circuit stays open (default: 30s)
}

type circuitBreakerProvider struct {
	provider llm.Provider
	cb       *gobreaker.CircuitBreaker
	name     string
}

func NewCircuitBreaker(provider llm.Provider, config CircuitBreakerConfig) llm.Provider {
	settings := gobreaker.Settings{
		Name:        provider.Name(),
		MaxRequests: uint32(config.SuccessThreshold),
		Interval:    config.Timeout,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests == 0 {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.TotalFailures >= uint32(config.FailureThreshold) && failureRatio >= 0.5
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			var pe *llm.ProviderError
			if errors.As(err, &pe) {
				return pe.StatusCode < 500
			}
			return false
		},
	}

	return &circuitBreakerProvider{
		provider: provider,
		cb:       gobreaker.NewCircuitBreaker(settings),
		name:     provider.Name(),
	}
}

func (c *circuitBreakerProvider) Name() string {
	return c.name
}

func (c *circuitBreakerProvider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	result, err := c.cb.Execute(func() (interface{}, error) {
		return c.provider.Chat(ctx, req)
	})

	if err != nil {
		logger.Log.Warn().
			Str("provider", c.name).
			Err(err).
			Msg("Circuit breaker error")
		return nil, err
	}

	return result.(*llm.ChatResponse), nil
}

func (c *circuitBreakerProvider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	_, err := c.cb.Execute(func() (interface{}, error) {
		err := c.provider.ChatStream(ctx, req, callback)
		return nil, err
	})

	if err != nil {
		logger.Log.Warn().
			Str("provider", c.name).
			Err(err).
			Msg("Circuit breaker stream error")
		return err
	}

	return nil
}

type CircuitBreakerManager struct {
	breakers      map[string]*gobreaker.CircuitBreaker
	defaultConfig CircuitBreakerConfig
}

func NewCircuitBreakerManager(config CircuitBreakerConfig) *CircuitBreakerManager {
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 5
	}
	if config.SuccessThreshold == 0 {
		config.SuccessThreshold = 3
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &CircuitBreakerManager{
		breakers:      make(map[string]*gobreaker.CircuitBreaker),
		defaultConfig: config,
	}
}

func (m *CircuitBreakerManager) GetBreaker(name string) *gobreaker.CircuitBreaker {
	return m.breakers[name]
}

func (m *CircuitBreakerManager) GetState(name string) string {
	if cb, ok := m.breakers[name]; ok {
		return cb.State().String()
	}
	return "unknown"
}

func (m *CircuitBreakerManager) WrapProvider(provider llm.Provider) llm.Provider {
	settings := gobreaker.Settings{
		Name:        provider.Name(),
		MaxRequests: uint32(m.defaultConfig.SuccessThreshold),
		Interval:    m.defaultConfig.Timeout,
		Timeout:     m.defaultConfig.Timeout,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			if counts.Requests == 0 {
				return false
			}
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.TotalFailures >= uint32(m.defaultConfig.FailureThreshold) && failureRatio >= 0.5
		},
		IsSuccessful: func(err error) bool {
			if err == nil {
				return true
			}
			var pe *llm.ProviderError
			if errors.As(err, &pe) {
				return pe.StatusCode < 500
			}
			return false
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)
	m.breakers[provider.Name()] = cb

	return &circuitBreakerProvider{
		provider: provider,
		cb:       cb,
		name:     provider.Name(),
	}
}

var (
	ErrCircuitOpen = errors.New("circuit breaker is open")
)

func IsCircuitOpen(err error) bool {
	return errors.Is(err, gobreaker.ErrOpenState)
}

func GetCircuitBreakerErrorMessage(err error) string {
	if IsCircuitOpen(err) {
		return "Service temporarily unavailable - provider circuit breaker is open"
	}
	return err.Error()
}
