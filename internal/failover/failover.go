package failover

import (
	"context"
	"strings"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/platform/logger"
)

type ProviderWithConfig struct {
	Provider        llm.Provider
	EnableRetries   bool
	MaxRetries      int
	RetryMultiplier int
}

type FallbackConfig struct {
	EnableRetries bool
	MaxRetries    int
}

var DefaultConfig = FallbackConfig{
	EnableRetries: false,
	MaxRetries:    0,
}

type failoverProvider struct {
	providers []ProviderWithConfig
	name      string
}

func ParseModelWithFallbacks(qualifiedModel string) []string {
	if !strings.Contains(qualifiedModel, "|") {
		return []string{qualifiedModel}
	}
	return strings.Split(qualifiedModel, "|")
}

func NewFailoverProvider(providers []ProviderWithConfig) llm.Provider {
	names := make([]string, len(providers))
	for i, p := range providers {
		names[i] = p.Provider.Name()
	}

	return &failoverProvider{
		providers: providers,
		name:      "failover(" + strings.Join(names, "->") + ")",
	}
}

func (f *failoverProvider) Name() string {
	return f.name
}

func (f *failoverProvider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	var lastErr error

	for i, p := range f.providers {
		logger.Log.Info().
			Str("provider", p.Provider.Name()).
			Int("fallback_index", i).
			Msg("Attempting provider")

		resp, err := p.Provider.Chat(ctx, req)
		if err == nil {
			if i > 0 {
				logger.Log.Info().
					Str("provider", p.Provider.Name()).
					Str("fallback_chain", f.name).
					Msg("Fallback succeeded")
			}
			return resp, nil
		}

		lastErr = err
		logger.Log.Warn().
			Str("provider", p.Provider.Name()).
			Err(err).
			Int("fallback_index", i).
			Msg("Provider failed, trying next fallback")
	}

	logger.Log.Error().
		Str("fallback_chain", f.name).
		Err(lastErr).
		Msg("All fallback providers failed")

	return nil, lastErr
}

func (f *failoverProvider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	var lastErr error

	for i, p := range f.providers {
		logger.Log.Info().
			Str("provider", p.Provider.Name()).
			Int("fallback_index", i).
			Msg("Attempting streaming provider")

		err := p.Provider.ChatStream(ctx, req, func(chunk *llm.StreamChunk) error {
			if i > 0 {
				logger.Log.Info().
					Str("provider", p.Provider.Name()).
					Msg("Streaming fallback succeeded")
			}
			return callback(chunk)
		})

		if err == nil {
			return nil
		}

		lastErr = err
		logger.Log.Warn().
			Str("provider", p.Provider.Name()).
			Err(err).
			Int("fallback_index", i).
			Msg("Streaming provider failed, trying next fallback")
	}

	logger.Log.Error().
		Str("fallback_chain", f.name).
		Err(lastErr).
		Msg("All streaming fallback providers failed")

	return lastErr
}
