package providers

import (
	"fmt"
	"strings"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/azure"
	"github.com/atozi-ai/gateway/internal/providers/openai"
	"github.com/atozi-ai/gateway/internal/providers/xai"
	"github.com/atozi-ai/gateway/internal/providers/zai"
)

// Get parses a qualified model name ("provider/model") and returns the
// matching provider plus the bare model name to send upstream.
func Get(qualifiedModel string, apiKey string, endpoint string) (llm.Provider, string, error) {
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
	case "azure":
		if endpoint == "" {
			return nil, "", &llm.ProviderError{
				StatusCode: 400,
				Message:    "azure provider requires an endpoint to be provided (e.g., https://<YOUR_RESOURCE_NAME>.openai.azure.com/openai/deployments/<YOUR_DEPLOYMENT_NAME>/chat/completions?api-version=<API_VERSION>)",
				Type:       "invalid_request_error",
				Code:       "missing_endpoint",
			}
		}
		return azure.New(apiKey, endpoint), model, nil
	case "xai":
		return xai.New(apiKey), model, nil
	case "zai":
		return zai.New(apiKey), model, nil
	default:
		return nil, "", &llm.ProviderError{
			StatusCode: 400,
			Message:    fmt.Sprintf("unknown provider: %q", providerName),
			Type:       "invalid_request_error",
			Code:       "unknown_provider",
		}
	}
}
