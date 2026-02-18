package providers

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/atozi-ai/gateway/internal/circuitbreaker"
	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/azure"
	"github.com/atozi-ai/gateway/internal/providers/openai"
	"github.com/atozi-ai/gateway/internal/providers/xai"
	"github.com/atozi-ai/gateway/internal/providers/zai"
)

type ProviderManager struct {
	mu        sync.RWMutex
	providers map[string]llm.Provider
	cbManager *circuitbreaker.CircuitBreakerManager
}

var (
	defaultManager *ProviderManager
	managerOnce    sync.Once
)

func GetProviderManager() *ProviderManager {
	managerOnce.Do(func() {
		defaultManager = &ProviderManager{
			providers: make(map[string]llm.Provider),
			cbManager: circuitbreaker.NewCircuitBreakerManager(circuitbreaker.CircuitBreakerConfig{
				FailureThreshold: 5,
				SuccessThreshold: 3,
				Timeout:          30 * time.Second,
			}),
		}
	})
	return defaultManager
}

func (m *ProviderManager) Get(qualifiedModel string, apiKey string, endpoint string) (llm.Provider, string, error) {
	providerName, model, ok := strings.Cut(qualifiedModel, "/")
	if !ok {
		return nil, "", &llm.ProviderError{
			StatusCode: 400,
			Message:    fmt.Sprintf("model must be in provider/model format, got %q", qualifiedModel),
			Type:       "invalid_request_error",
			Code:       "invalid_model_format",
		}
	}

	provider, err := m.getProvider(providerName, apiKey, endpoint)
	if err != nil {
		return nil, "", err
	}

	return provider, model, nil
}

func (m *ProviderManager) getProvider(name string, apiKey string, endpoint string) (llm.Provider, error) {
	cacheKey := fmt.Sprintf("%s:%s:%s", name, apiKey, endpoint)

	m.mu.RLock()
	if provider, exists := m.providers[cacheKey]; exists {
		m.mu.RUnlock()
		return provider, nil
	}
	m.mu.RUnlock()

	m.mu.Lock()
	defer m.mu.Unlock()

	if provider, exists := m.providers[cacheKey]; exists {
		return provider, nil
	}

	var provider llm.Provider
	switch name {
	case "openai":
		provider = openai.New(apiKey)
	case "azure":
		if endpoint == "" {
			return nil, &llm.ProviderError{
				StatusCode: 400,
				Message:    "azure provider requires an endpoint to be provided",
				Type:       "invalid_request_error",
				Code:       "missing_endpoint",
			}
		}
		provider = azure.New(apiKey, endpoint)
	case "xai":
		provider = xai.New(apiKey)
	case "zai":
		provider = zai.New(apiKey)
	default:
		return nil, &llm.ProviderError{
			StatusCode: 400,
			Message:    fmt.Sprintf("unknown provider: %q", name),
			Type:       "invalid_request_error",
			Code:       "unknown_provider",
		}
	}

	wrappedProvider := m.cbManager.WrapProvider(provider)
	m.providers[cacheKey] = wrappedProvider

	return wrappedProvider, nil
}

func (m *ProviderManager) GetCircuitBreakerState(name string) string {
	return m.cbManager.GetState(name)
}

// Get parses a qualified model name ("provider/model") and returns the
// matching provider plus the bare model name to send upstream.
// Uses the default global ProviderManager with circuit breaker.
func Get(qualifiedModel string, apiKey string, endpoint string) (llm.Provider, string, error) {
	return GetProviderManager().Get(qualifiedModel, apiKey, endpoint)
}
