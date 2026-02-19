package providers

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/atozi-ai/gateway/internal/circuitbreaker"
	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/failover"
	"github.com/atozi-ai/gateway/internal/platform/logger"
	"github.com/atozi-ai/gateway/internal/providers/ai21"
	"github.com/atozi-ai/gateway/internal/providers/anthropic"
	"github.com/atozi-ai/gateway/internal/providers/anyscale"
	"github.com/atozi-ai/gateway/internal/providers/azure"
	"github.com/atozi-ai/gateway/internal/providers/baseten"
	"github.com/atozi-ai/gateway/internal/providers/cerebras"
	"github.com/atozi-ai/gateway/internal/providers/cloudflare"
	"github.com/atozi-ai/gateway/internal/providers/cohere"
	"github.com/atozi-ai/gateway/internal/providers/deepinfra"
	"github.com/atozi-ai/gateway/internal/providers/deepseek"
	"github.com/atozi-ai/gateway/internal/providers/fireworks"
	"github.com/atozi-ai/gateway/internal/providers/friendli"
	"github.com/atozi-ai/gateway/internal/providers/gemini"
	"github.com/atozi-ai/gateway/internal/providers/groq"
	"github.com/atozi-ai/gateway/internal/providers/hyperbolic"
	"github.com/atozi-ai/gateway/internal/providers/minimax"
	"github.com/atozi-ai/gateway/internal/providers/mistral"
	"github.com/atozi-ai/gateway/internal/providers/moonshot"
	"github.com/atozi-ai/gateway/internal/providers/nebius"
	"github.com/atozi-ai/gateway/internal/providers/novita"
	"github.com/atozi-ai/gateway/internal/providers/nvidia"
	"github.com/atozi-ai/gateway/internal/providers/ollama"
	"github.com/atozi-ai/gateway/internal/providers/openai"
	"github.com/atozi-ai/gateway/internal/providers/perplexity"
	"github.com/atozi-ai/gateway/internal/providers/replicate"
	"github.com/atozi-ai/gateway/internal/providers/sambanova"
	"github.com/atozi-ai/gateway/internal/providers/siliconflow"
	"github.com/atozi-ai/gateway/internal/providers/together"
	"github.com/atozi-ai/gateway/internal/providers/upstage"
	"github.com/atozi-ai/gateway/internal/providers/xai"
	"github.com/atozi-ai/gateway/internal/providers/zai"
	"github.com/atozi-ai/gateway/internal/retry"
)

type ProviderManager struct {
	mu                      sync.RWMutex
	providers               map[string]llm.Provider
	cbManager               *circuitbreaker.CircuitBreakerManager
	enableRetryWithFallback bool
}

var (
	defaultManager *ProviderManager
	managerOnce    sync.Once
)

func GetProviderManager() *ProviderManager {
	managerOnce.Do(func() {
		retryWithFallback := os.Getenv("RETRY_WITH_FALLBACK")
		enableRetryWithFallback := retryWithFallback == "true" || retryWithFallback == "1"

		defaultManager = &ProviderManager{
			providers:               make(map[string]llm.Provider),
			enableRetryWithFallback: enableRetryWithFallback,
			cbManager: circuitbreaker.NewCircuitBreakerManager(circuitbreaker.CircuitBreakerConfig{
				FailureThreshold: 5,
				SuccessThreshold: 3,
				Timeout:          30 * time.Second,
			}),
		}

		logger.Log.Info().
			Bool("enable_retry_with_fallback", enableRetryWithFallback).
			Msg("Provider manager initialized")
	})
	return defaultManager
}

func (m *ProviderManager) Get(qualifiedModel string, apiKey string, endpoint string) (llm.Provider, string, error) {
	models := failover.ParseModelWithFallbacks(qualifiedModel)

	if len(models) == 1 {
		providerName, model, ok := strings.Cut(models[0], "/")
		if !ok {
			return nil, "", &llm.ProviderError{
				StatusCode: 400,
				Message:    fmt.Sprintf("model must be in provider/model format, got %q", qualifiedModel),
				Type:       "invalid_request_error",
				Code:       "invalid_model_format",
			}
		}

		provider, err := m.getProvider(providerName, endpoint, true)
		if err != nil {
			return nil, "", err
		}

		return provider, model, nil
	}

	var providersWithConfig []failover.ProviderWithConfig
	var finalModel string
	enableRetries := m.enableRetryWithFallback

	for i, modelSpec := range models {
		providerName, model, ok := strings.Cut(modelSpec, "/")
		if !ok {
			return nil, "", &llm.ProviderError{
				StatusCode: 400,
				Message:    fmt.Sprintf("model must be in provider/model format, got %q", modelSpec),
				Type:       "invalid_request_error",
				Code:       "invalid_model_format",
			}
		}

		if i == 0 {
			finalModel = model
		}

		provider, err := m.getProvider(providerName, endpoint, enableRetries)
		if err != nil {
			logger.Log.Warn().
				Str("model_spec", modelSpec).
				Err(err).
				Msg("Failed to create provider for fallback")
			continue
		}

		providersWithConfig = append(providersWithConfig, failover.ProviderWithConfig{
			Provider:      provider,
			EnableRetries: enableRetries,
		})
	}

	if len(providersWithConfig) == 0 {
		return nil, "", &llm.ProviderError{
			StatusCode: 400,
			Message:    "no valid fallback providers available",
			Type:       "invalid_request_error",
			Code:       "no_providers",
		}
	}

	failoverProvider := failover.NewFailoverProvider(providersWithConfig)
	return failoverProvider, finalModel, nil
}

func (m *ProviderManager) getProvider(name string, endpoint string, enableRetry bool) (llm.Provider, error) {
	cacheKey := fmt.Sprintf("%s:%s:%v", name, endpoint, enableRetry)

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

	var baseProvider llm.Provider
	switch name {
	case "openai":
		baseProvider = openai.New()
	case "azure":
		if endpoint == "" {
			return nil, &llm.ProviderError{
				StatusCode: 400,
				Message:    "azure provider requires an endpoint to be provided",
				Type:       "invalid_request_error",
				Code:       "missing_endpoint",
			}
		}
		baseProvider = azure.New("", endpoint)
	case "ai21":
		baseProvider = ai21.New()
	case "baseten":
		baseProvider = baseten.New()
	case "anyscale":
		baseProvider = anyscale.New()
	case "cerebras":
		baseProvider = cerebras.New()
	case "cloudflare":
		baseProvider = cloudflare.New()
	case "xai":
		baseProvider = xai.New()
	case "zai":
		baseProvider = zai.New()
	case "anthropic":
		baseProvider = anthropic.New()
	case "gemini":
		baseProvider = gemini.New()
	case "groq":
		baseProvider = groq.New()
	case "hyperbolic":
		baseProvider = hyperbolic.New()
	case "minimax":
		baseProvider = minimax.New()
	case "deepinfra":
		baseProvider = deepinfra.New()
	case "deepseek":
		baseProvider = deepseek.New()
	case "mistral":
		baseProvider = mistral.New()
	case "moonshot":
		baseProvider = moonshot.New()
	case "nebius":
		baseProvider = nebius.New()
	case "nvidia":
		baseProvider = nvidia.New()
	case "together":
		baseProvider = together.New()
	case "upstage":
		baseProvider = upstage.New()
	case "fireworks":
		baseProvider = fireworks.New()
	case "friendli":
		baseProvider = friendli.New()
	case "perplexity":
		baseProvider = perplexity.New()
	case "replicate":
		baseProvider = replicate.New()
	case "sambanova":
		baseProvider = sambanova.New()
	case "cohere":
		baseProvider = cohere.New()
	case "novita":
		baseProvider = novita.New()
	case "ollama":
		baseProvider = ollama.New()
	case "siliconflow":
		baseProvider = siliconflow.New()
	default:
		return nil, &llm.ProviderError{
			StatusCode: 400,
			Message:    fmt.Sprintf("unknown provider: %q", name),
			Type:       "invalid_request_error",
			Code:       "unknown_provider",
		}
	}

	wrappedProvider := m.cbManager.WrapProvider(baseProvider)

	if enableRetry {
		wrappedProvider = retry.NewRetryableProvider(wrappedProvider, retry.Config{
			MaxRetries:   3,
			InitialDelay: 500 * time.Millisecond,
			MaxDelay:     10 * time.Second,
			Multiplier:   2.0,
		})
	}

	m.providers[cacheKey] = wrappedProvider

	return wrappedProvider, nil
}

func (m *ProviderManager) GetCircuitBreakerState(name string) string {
	return m.cbManager.GetState(name)
}

func Get(qualifiedModel string, apiKey string, endpoint string) (llm.Provider, string, error) {
	return GetProviderManager().Get(qualifiedModel, apiKey, endpoint)
}
