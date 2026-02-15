package providers

import (
	"strings"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai"
)

func Get(model string) llm.Provider {
	m := strings.ToLower(model)

	switch {
	case strings.Contains(m, "openai"):
		return openai.New()
	default:
		return openai.New()
	}
}
