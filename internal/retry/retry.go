package retry

import (
	"context"
	"errors"
	"math"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
)

type Config struct {
	MaxRetries     int
	InitialDelay   time.Duration
	MaxDelay       time.Duration
	Multiplier     float64
	RetryableCodes []int
}

var DefaultConfig = Config{
	MaxRetries:   3,
	InitialDelay: 500 * time.Millisecond,
	MaxDelay:     10 * time.Second,
	Multiplier:   2.0,
	RetryableCodes: []int{
		429, // Rate limit
		500, // Internal server error
		502, // Bad gateway
		503, // Service unavailable
		504, // Gateway timeout
	},
}

func isRetryable(err error, retryableCodes []int) bool {
	var pe *llm.ProviderError
	if errors.As(err, &pe) {
		for _, code := range retryableCodes {
			if pe.StatusCode == code {
				return true
			}
		}
	}
	return false
}

func calculateDelay(attempt int, config Config) time.Duration {
	delay := float64(config.InitialDelay) * math.Pow(config.Multiplier, float64(attempt))
	if delay > float64(config.MaxDelay) {
		return config.MaxDelay
	}
	return time.Duration(delay)
}

type retryableProvider struct {
	provider llm.Provider
	config   Config
}

func NewRetryableProvider(provider llm.Provider, config Config) llm.Provider {
	if config.MaxRetries == 0 {
		config.MaxRetries = DefaultConfig.MaxRetries
	}
	if config.InitialDelay == 0 {
		config.InitialDelay = DefaultConfig.InitialDelay
	}
	if config.MaxDelay == 0 {
		config.MaxDelay = DefaultConfig.MaxDelay
	}
	if config.Multiplier == 0 {
		config.Multiplier = DefaultConfig.Multiplier
	}
	if len(config.RetryableCodes) == 0 {
		config.RetryableCodes = DefaultConfig.RetryableCodes
	}

	return &retryableProvider{
		provider: provider,
		config:   config,
	}
}

func (r *retryableProvider) Name() string {
	return r.provider.Name()
}

func (r *retryableProvider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := calculateDelay(attempt-1, r.config)
			logger.Log.Info().
				Str("provider", r.provider.Name()).
				Int("attempt", attempt).
				Dur("delay", delay).
				Msg("Retrying request")

			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		resp, err := r.provider.Chat(ctx, req)
		if err == nil {
			if attempt > 0 {
				logger.Log.Info().
					Str("provider", r.provider.Name()).
					Int("attempts", attempt+1).
					Msg("Request succeeded after retry")
			}
			return resp, nil
		}

		lastErr = err

		if !isRetryable(err, r.config.RetryableCodes) {
			logger.Log.Warn().
				Str("provider", r.provider.Name()).
				Err(err).
				Int("status_code", getStatusCode(err)).
				Msg("Non-retryable error, not retrying")
			return nil, err
		}

		logger.Log.Warn().
			Str("provider", r.provider.Name()).
			Err(err).
			Int("attempt", attempt+1).
			Int("max_retries", r.config.MaxRetries).
			Int("status_code", getStatusCode(err)).
			Msg("Retryable error, will retry")
	}

	logger.Log.Error().
		Str("provider", r.provider.Name()).
		Err(lastErr).
		Int("max_retries", r.config.MaxRetries).
		Msg("All retry attempts exhausted")

	return nil, lastErr
}

func (r *retryableProvider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	var lastErr error

	for attempt := 0; attempt <= r.config.MaxRetries; attempt++ {
		if attempt > 0 {
			delay := calculateDelay(attempt-1, r.config)
			logger.Log.Info().
				Str("provider", r.provider.Name()).
				Int("attempt", attempt).
				Dur("delay", delay).
				Msg("Retrying streaming request")

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(delay):
			}
		}

		err := r.provider.ChatStream(ctx, req, func(chunk *llm.StreamChunk) error {
			if attempt > 0 && chunk.Choices != nil && len(chunk.Choices) > 0 {
				logger.Log.Info().
					Str("provider", r.provider.Name()).
					Msg("Streaming succeeded after retry")
			}
			return callback(chunk)
		})

		if err == nil {
			if attempt > 0 {
				logger.Log.Info().
					Str("provider", r.provider.Name()).
					Int("attempts", attempt+1).
					Msg("Streaming request succeeded after retry")
			}
			return nil
		}

		lastErr = err

		if !isRetryable(err, r.config.RetryableCodes) {
			logger.Log.Warn().
				Str("provider", r.provider.Name()).
				Err(err).
				Int("status_code", getStatusCode(err)).
				Msg("Non-retryable error for streaming, not retrying")
			return err
		}

		logger.Log.Warn().
			Str("provider", r.provider.Name()).
			Err(err).
			Int("attempt", attempt+1).
			Int("max_retries", r.config.MaxRetries).
			Int("status_code", getStatusCode(err)).
			Msg("Retryable error for streaming, will retry")
	}

	logger.Log.Error().
		Str("provider", r.provider.Name()).
		Err(lastErr).
		Int("max_retries", r.config.MaxRetries).
		Msg("All retry attempts exhausted for streaming")

	return lastErr
}

func getStatusCode(err error) int {
	var pe *llm.ProviderError
	if errors.As(err, &pe) {
		return pe.StatusCode
	}
	return 0
}
