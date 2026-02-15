package providers

import (
	"strings"
	"sync"

	"github.com/atozi-ai/gateway/internal/domain/llm"
	"github.com/atozi-ai/gateway/internal/providers/openai"
)

var (
	// Cache for provider instances (singleton pattern)
	providerCache = make(map[string]llm.Provider)
	cacheMu       sync.RWMutex
)

// Get returns a cached provider instance for the given model.
// Provider instances are created once and reused for all requests.
func Get(model string) llm.Provider {
	m := strings.ToLower(model)

	// Determine provider type based on model name
	var providerType string
	switch {
	case strings.Contains(m, "openai"):
		providerType = "openai"
	default:
		providerType = "openai"
	}

	// Check cache first (read lock)
	cacheMu.RLock()
	if provider, exists := providerCache[providerType]; exists {
		cacheMu.RUnlock()
		return provider
	}
	cacheMu.RUnlock()

	// Create new provider instance (write lock)
	cacheMu.Lock()
	defer cacheMu.Unlock()

	// Double-check pattern: another goroutine might have created it
	if provider, exists := providerCache[providerType]; exists {
		return provider
	}

	// Create and cache the provider
	var provider llm.Provider
	switch providerType {
	case "openai":
		provider = openai.New()
	default:
		provider = openai.New()
	}

	providerCache[providerType] = provider
	return provider
}
