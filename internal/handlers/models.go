package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ModelInfo struct {
	ID         string `json:"id"`
	Object     string `json:"object"`
	OwnedBy    string `json:"owned_by"`
	Provider   string `json:"provider"`
	Name       string `json:"name"`
	ContextLen int    `json:"context_len"`
}

type ModelsListResponse struct {
	Object string      `json:"object"`
	Data   []ModelInfo `json:"data"`
}

type ModelsHandler struct{}

func NewModelsHandler() *ModelsHandler {
	return &ModelsHandler{}
}

func (h *ModelsHandler) ListModels(w http.ResponseWriter, r *http.Request) {
	models := []ModelInfo{
		// OpenAI Models - Latest 2025-2026
		{ID: "openai/gpt-5.2", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-5.2", ContextLen: 256000},
		{ID: "openai/gpt-4.5", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.5", ContextLen: 128000},
		{ID: "openai/gpt-4.1", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.1", ContextLen: 200000},
		{ID: "openai/gpt-4.1-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.1 Mini", ContextLen: 200000},
		{ID: "openai/gpt-4o", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o", ContextLen: 128000},
		{ID: "openai/gpt-4o-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o Mini", ContextLen: 128000},
		{ID: "openai/o1", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1", ContextLen: 200000},
		{ID: "openai/o1-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1 Mini", ContextLen: 128000},
		{ID: "openai/o3", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o3", ContextLen: 200000},
		{ID: "openai/o3-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o3 Mini", ContextLen: 200000},
		{ID: "openai/o4-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o4 Mini", ContextLen: 200000},
		{ID: "openai/gpt-4-turbo", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4 Turbo", ContextLen: 128000},
		{ID: "openai/gpt-4", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4", ContextLen: 8192},
		{ID: "openai/gpt-3.5-turbo", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-3.5 Turbo", ContextLen: 16385},
		{ID: "openai/gpt-oss-120b", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-OSS 120B", ContextLen: 128000},
		{ID: "openai/gpt-oss-20b", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-OSS 20B", ContextLen: 128000},

		// xAI Grok Models - Latest 2025-2026
		{ID: "xai/grok-4", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4", ContextLen: 256000},
		{ID: "xai/grok-4-fast", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4 Fast", ContextLen: 2000000},
		{ID: "xai/grok-4.1", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4.1", ContextLen: 256000},
		{ID: "xai/grok-4.1-fast-reasoning", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4.1 Fast Reasoning", ContextLen: 2000000},
		{ID: "xai/grok-3", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 3", ContextLen: 131072},
		{ID: "xai/grok-3-mini", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 3 Mini", ContextLen: 131072},
		{ID: "xai/grok-2", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 2", ContextLen: 131072},
		{ID: "xai/grok-2-vision", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 2 Vision", ContextLen: 32768},

		// ZAI Models
		{ID: "zai/zai-core", Object: "model", OwnedBy: "zai", Provider: "zai", Name: "ZAI Core", ContextLen: 128000},
		{ID: "zai/zai-core-2025-01-20", Object: "model", OwnedBy: "zai", Provider: "zai", Name: "ZAI Core 2025-01-20", ContextLen: 128000},

		// Azure OpenAI Models
		{ID: "azure/gpt-4o", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4o", ContextLen: 128000},
		{ID: "azure/gpt-4o-mini", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4o Mini", ContextLen: 128000},
		{ID: "azure/gpt-4-turbo", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4 Turbo", ContextLen: 128000},
		{ID: "azure/gpt-4", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4", ContextLen: 8192},
		{ID: "azure/gpt-35-turbo", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-3.5 Turbo", ContextLen: 16385},

		// Anthropic Claude Models - Latest 2025-2026
		{ID: "anthropic/claude-opus-4-6", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Opus 4.6", ContextLen: 200000},
		{ID: "anthropic/claude-opus-4-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Opus 4.5", ContextLen: 200000},
		{ID: "anthropic/claude-opus-4-1", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Opus 4.1", ContextLen: 200000},
		{ID: "anthropic/claude-sonnet-4-6", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 4.6", ContextLen: 200000},
		{ID: "anthropic/claude-sonnet-4-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 4.5", ContextLen: 200000},
		{ID: "anthropic/claude-sonnet-4", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 4", ContextLen: 200000},
		{ID: "anthropic/claude-sonnet-3-7", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 3.7", ContextLen: 200000},
		{ID: "anthropic/claude-haiku-4-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Haiku 4.5", ContextLen: 200000},
		{ID: "anthropic/claude-haiku-3-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Haiku 3.5", ContextLen: 200000},

		// Google Gemini Models - Latest 2025-2026
		{ID: "gemini/gemini-3-pro-preview", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 3 Pro Preview", ContextLen: 1048576},
		{ID: "gemini/gemini-3-flash-preview", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 3 Flash Preview", ContextLen: 1048576},
		{ID: "gemini/gemini-2.5-flash", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.5 Flash", ContextLen: 1000000},
		{ID: "gemini/gemini-2.5-pro", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.5 Pro", ContextLen: 1000000},
		{ID: "gemini/gemini-2.0-flash", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.0 Flash", ContextLen: 1000000},
		{ID: "gemini/gemini-2.0-pro", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.0 Pro", ContextLen: 2000000},
		{ID: "gemini/gemini-2.0-flash-lite", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.0 Flash Lite", ContextLen: 1000000},
		{ID: "gemini/gemini-1.5-flash", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 1.5 Flash", ContextLen: 1000000},
		{ID: "gemini/gemini-1.5-pro", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 1.5 Pro", ContextLen: 2000000},

		// DeepSeek Models - 2025
		{ID: "deepseek/deepseek-chat", Object: "model", OwnedBy: "deepseek", Provider: "deepseek", Name: "DeepSeek V3 Chat", ContextLen: 128000},
		{ID: "deepseek/deepseek-reasoner", Object: "model", OwnedBy: "deepseek", Provider: "deepseek", Name: "DeepSeek R1 Reasoner", ContextLen: 128000},

		// Mistral AI Models - Latest 2025-2026
		{ID: "mistral/mistral-large-3", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Mistral Large 3", ContextLen: 256000},
		{ID: "mistral/mistral-medium-3-1", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Mistral Medium 3.1", ContextLen: 128000},
		{ID: "mistral/mistral-small-3-2", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Mistral Small 3.2", ContextLen: 32000},
		{ID: "mistral/devstral-2", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Devstral 2", ContextLen: 128000},
		{ID: "mistral/magistral-1-2", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Magistral 1.2", ContextLen: 128000},
		{ID: "mistral/codestral-2508", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Codestral 25.08", ContextLen: 128000},
		{ID: "mistral/ministral-14b", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Ministral 14B", ContextLen: 128000},
		{ID: "mistral/ministral-8b", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Ministral 8B", ContextLen: 128000},
		{ID: "mistral/ministral-3b", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Ministral 3B", ContextLen: 128000},
		{ID: "mistral/pixtral-large", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Pixtral Large", ContextLen: 128000},
		{ID: "mistral/voxtral", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Voxtral", ContextLen: 128000},
		{ID: "mistral/nemo", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Nemo", ContextLen: 128000},

		// Groq Models - Latest 2025-2026
		{ID: "groq/llama-3.3-70b-versatile", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "groq/llama-3.1-8b-instant", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 3.1 8B", ContextLen: 131072},
		{ID: "groq/llama-4-maverick", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 4 Maverick", ContextLen: 131072},
		{ID: "groq/llama-4-scout", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 4 Scout", ContextLen: 131072},
		{ID: "groq/openai-gpt-oss-120b", Object: "model", OwnedBy: "openai", Provider: "groq", Name: "GPT-OSS 120B", ContextLen: 131072},
		{ID: "groq/openai-gpt-oss-20b", Object: "model", OwnedBy: "openai", Provider: "groq", Name: "GPT-OSS 20B", ContextLen: 131072},
		{ID: "groq/qwen-2.5-32b", Object: "model", OwnedBy: "alibaba", Provider: "groq", Name: "Qwen 2.5 32B", ContextLen: 131072},
		{ID: "groq/qwen-2.5-coder-32b", Object: "model", OwnedBy: "alibaba", Provider: "groq", Name: "Qwen 2.5 Coder 32B", ContextLen: 131072},
		{ID: "groq/qwen-qwq-32b", Object: "model", OwnedBy: "alibaba", Provider: "groq", Name: "Qwen QwQ 32B", ContextLen: 131072},
		{ID: "groq/deepseek-r1-distill-qwen-32b", Object: "model", OwnedBy: "deepseek", Provider: "groq", Name: "DeepSeek R1 Distill Qwen 32B", ContextLen: 131072},
		{ID: "groq/gemma2-9b-it", Object: "model", OwnedBy: "google", Provider: "groq", Name: "Gemma 2 9B", ContextLen: 8192},
		{ID: "groq/mistral-saba-24b", Object: "model", OwnedBy: "mistral", Provider: "groq", Name: "Mistral Saba 24B", ContextLen: 32000},

		// Together AI Models - Latest 2025-2026
		{ID: "together/llama-3.3-70b-instruct-turbo", Object: "model", OwnedBy: "meta", Provider: "together", Name: "Llama 3.3 70B Turbo", ContextLen: 131072},
		{ID: "together/llama-3.1-405b-instruct-turbo", Object: "model", OwnedBy: "meta", Provider: "together", Name: "Llama 3.1 405B Turbo", ContextLen: 131072},
		{ID: "together/llama-3.1-70b-instruct-turbo", Object: "model", OwnedBy: "meta", Provider: "together", Name: "Llama 3.1 70B Turbo", ContextLen: 131072},
		{ID: "together/llama-3.1-8b-instruct-turbo", Object: "model", OwnedBy: "meta", Provider: "together", Name: "Llama 3.1 8B Turbo", ContextLen: 131072},
		{ID: "together/deepseek-v3.1", Object: "model", OwnedBy: "deepseek", Provider: "together", Name: "DeepSeek V3.1", ContextLen: 128000},
		{ID: "together/qwen3.5-397b-a17b", Object: "model", OwnedBy: "qwen", Provider: "together", Name: "Qwen 3.5 397B A17B", ContextLen: 262144},
		{ID: "together/qwen2.5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "together", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},
		{ID: "together/qwen2.5-32b-instruct", Object: "model", OwnedBy: "qwen", Provider: "together", Name: "Qwen 2.5 32B Instruct", ContextLen: 131072},
		{ID: "together/mistral-large-3", Object: "model", OwnedBy: "mistral", Provider: "together", Name: "Mistral Large 3", ContextLen: 131072},
		{ID: "together/mistral-medium-3-1", Object: "model", OwnedBy: "mistral", Provider: "together", Name: "Mistral Medium 3.1", ContextLen: 131072},
		{ID: "together/moonshotai-kimi-k2-instruct", Object: "model", OwnedBy: "moonshot", Provider: "together", Name: "Kimi K2 Instruct", ContextLen: 262144},
		{ID: "together/minimax-m2.5", Object: "model", OwnedBy: "minimax", Provider: "together", Name: "MiniMax M2.5", ContextLen: 228700},

		// Fireworks AI Models - Latest 2025-2026
		{ID: "fireworks/llama-v3p3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "fireworks", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "fireworks/llama-v3p1-405b-instruct", Object: "model", OwnedBy: "meta", Provider: "fireworks", Name: "Llama 3.1 405B Instruct", ContextLen: 131072},
		{ID: "fireworks/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "fireworks", Name: "DeepSeek V3", ContextLen: 128000},
		{ID: "fireworks/qwen2p5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "fireworks", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},
		{ID: "fireworks/qwen3-coder-480b-a35b", Object: "model", OwnedBy: "qwen", Provider: "fireworks", Name: "Qwen3 Coder 480B A35B", ContextLen: 262144},
		{ID: "fireworks/openai-gpt-oss-120b", Object: "model", OwnedBy: "openai", Provider: "fireworks", Name: "GPT-OSS 120B", ContextLen: 131072},
		{ID: "fireworks/openai-gpt-oss-20b", Object: "model", OwnedBy: "openai", Provider: "fireworks", Name: "GPT-OSS 20B", ContextLen: 131072},
		{ID: "fireworks/glm-4.5", Object: "model", OwnedBy: "zhipu", Provider: "fireworks", Name: "GLM 4.5", ContextLen: 131072},
		{ID: "fireworks/kimi-k2-instruct", Object: "model", OwnedBy: "moonshot", Provider: "fireworks", Name: "Kimi K2 Instruct", ContextLen: 131072},
		{ID: "fireworks/mistral-large-3", Object: "model", OwnedBy: "mistral", Provider: "fireworks", Name: "Mistral Large 3", ContextLen: 131072},
		{ID: "fireworks/mixtral-8x22b-instruct", Object: "model", OwnedBy: "mistral", Provider: "fireworks", Name: "Mixtral 8x22B", ContextLen: 65536},
		{ID: "fireworks/gemma-2-9b-it", Object: "model", OwnedBy: "google", Provider: "fireworks", Name: "Gemma 2 9B", ContextLen: 8192},

		// Perplexity AI Models - Latest 2025-2026
		{ID: "perplexity/sonar", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar", ContextLen: 128000},
		{ID: "perplexity/sonar-pro", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar Pro", ContextLen: 200000},
		{ID: "perplexity/sonar-reasoning", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar Reasoning", ContextLen: 128000},
		{ID: "perplexity/sonar-reasoning-pro", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar Reasoning Pro", ContextLen: 128000},

		// Cohere Models - Latest 2025-2026
		{ID: "cohere/command-a", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Command A", ContextLen: 256000},
		{ID: "cohere/command-r-plus", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Command R+", ContextLen: 128000},
		{ID: "cohere/command-r", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Command R", ContextLen: 128000},
		{ID: "cohere/command-light", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Command Light", ContextLen: 4096},
		{ID: "cohere/embed-english-v3", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Embed English v3", ContextLen: 512},
		{ID: "cohere/embed-multilingual-v3", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Embed Multilingual v3", ContextLen: 512},
		{ID: "cohere/rerank-v3.5", Object: "model", OwnedBy: "cohere", Provider: "cohere", Name: "Rerank v3.5", ContextLen: 512},

		// Novita AI Models - Latest 2025-2026
		{ID: "novita/llama-3.3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "novita", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "novita/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "novita", Name: "DeepSeek V3", ContextLen: 163840},
		{ID: "novita/deepseek-v3-turbo", Object: "model", OwnedBy: "deepseek", Provider: "novita", Name: "DeepSeek V3 Turbo", ContextLen: 128000},
		{ID: "novita/deepseek-r1", Object: "model", OwnedBy: "deepseek", Provider: "novita", Name: "DeepSeek R1", ContextLen: 128000},
		{ID: "novita/qwen-2.5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "novita", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},
		{ID: "novita/qwen-2.5-coder-32b-instruct", Object: "model", OwnedBy: "qwen", Provider: "novita", Name: "Qwen 2.5 Coder 32B Instruct", ContextLen: 131072},
		{ID: "novita/mistral-nemo-instruct", Object: "model", OwnedBy: "mistral", Provider: "novita", Name: "Mistral Nemo Instruct", ContextLen: 128000},
		{ID: "novita/gemma-2-9b-it", Object: "model", OwnedBy: "google", Provider: "novita", Name: "Gemma 2 9B Instruct", ContextLen: 8192},

		// Hyperbolic Models - Latest 2025-2026
		{ID: "hyperbolic/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "hyperbolic", Name: "DeepSeek V3", ContextLen: 131072},
		{ID: "hyperbolic/deepseek-r1", Object: "model", OwnedBy: "deepseek", Provider: "hyperbolic", Name: "DeepSeek R1", ContextLen: 131072},
		{ID: "hyperbolic/qwen-2.5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "hyperbolic", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},
		{ID: "hyperbolic/llama-3.1-405b-base", Object: "model", OwnedBy: "meta", Provider: "hyperbolic", Name: "Llama 3.1 405B Base", ContextLen: 131072},
		{ID: "hyperbolic/gpt-oss-120b", Object: "model", OwnedBy: "openai", Provider: "hyperbolic", Name: "GPT-OSS 120B", ContextLen: 131069},

		// Upstage Models - Latest 2025-2026
		{ID: "upstage/solar-pro-3", Object: "model", OwnedBy: "upstage", Provider: "upstage", Name: "Solar Pro 3", ContextLen: 128000},
		{ID: "upstage/solar-pro-2", Object: "model", OwnedBy: "upstage", Provider: "upstage", Name: "Solar Pro 2", ContextLen: 128000},
		{ID: "upstage/solar-mini", Object: "model", OwnedBy: "upstage", Provider: "upstage", Name: "Solar Mini", ContextLen: 128000},
		{ID: "upstage/solar-open", Object: "model", OwnedBy: "upstage", Provider: "upstage", Name: "Solar Open", ContextLen: 128000},
	}

	response := ModelsListResponse{
		Object: "list",
		Data:   models,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ModelsHandler) RegisterRoutes(r chi.Router) {
	r.Get("/models", h.ListModels)
}
