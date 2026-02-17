package providers

import (
	"fmt"
	"strings"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai"
	"github.com/atozi-ai/gateway/internal/providers/xai"
)

// Get parses a qualified model name ("provider/model") and returns the
// matching provider plus the bare model name to send upstream.
func Get(qualifiedModel string, apiKey string) (llm.Provider, string, error) {
	providerName, model, ok := strings.Cut(qualifiedModel, "/")
	if !ok {
		return nil, "", &llm.ProviderError{
			StatusCode: 400,
			Message:    fmt.Sprintf("model must be in provider/model format, got %q", qualifiedModel),
			Type:       "invalid_request_error",
			Code:       "invalid_model_format",
		}
	}

	switch providerName {
	case "openai":
		return openai.New(apiKey), model, nil
	case "xai":
		return xai.New(apiKey), model, nil
	default:
		return nil, "", &llm.ProviderError{
			StatusCode: 400,
			Message:    fmt.Sprintf("unknown provider: %q", providerName),
			Type:       "invalid_request_error",
			Code:       "unknown_provider",
		}
	}
}
