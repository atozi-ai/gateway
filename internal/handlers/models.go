package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type ModelInfo struct {
	ID          string   `json:"id"`
	Object      string   `json:"object"`
	OwnedBy     string   `json:"owned_by"`
	Provider    string   `json:"provider"`
	Name        string   `json:"name"`
	ContextLen  int      `json:"context_len"`
	Description string   `json:"description,omitempty"`
	Category    []string `json:"category,omitempty"`
	IsFlagship  bool     `json:"is_flagship,omitempty"`
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
		// AWS Bedrock Models - Latest 2025-2026
		{ID: "bedrock/anthropic.claude-3-opus-20240229", Object: "model", OwnedBy: "anthropic", Provider: "bedrock", Name: "Claude 3 Opus", ContextLen: 200000, Description: "Anthropic's most capable model", Category: []string{"general", "reasoning", "coding"}, IsFlagship: true},
		{ID: "bedrock/anthropic.claude-3-sonnet-20240229", Object: "model", OwnedBy: "anthropic", Provider: "bedrock", Name: "Claude 3 Sonnet", ContextLen: 200000, Description: "Balanced performance model", Category: []string{"general", "coding"}},
		{ID: "bedrock/anthropic.claude-3-haiku-20240307", Object: "model", OwnedBy: "anthropic", Provider: "bedrock", Name: "Claude 3 Haiku", ContextLen: 200000, Description: "Fast efficient model", Category: []string{"general"}},
		{ID: "bedrock/anthropic.claude-sonnet-4-20250514", Object: "model", OwnedBy: "anthropic", Provider: "bedrock", Name: "Claude Sonnet 4", ContextLen: 200000, Description: "Latest balanced model", Category: []string{"general", "coding"}, IsFlagship: true},
		{ID: "bedrock/meta.llama3-3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "bedrock", Name: "Llama 3.3 70B Instruct", ContextLen: 131072, Description: "Meta's latest open model", Category: []string{"general"}, IsFlagship: true},
		{ID: "bedrock/meta.llama3-1-405b-instruct", Object: "model", OwnedBy: "meta", Provider: "bedrock", Name: "Llama 3.1 405B Instruct", ContextLen: 131072, Description: "Large Meta model", Category: []string{"general"}},
		{ID: "bedrock/meta.llama3-1-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "bedrock", Name: "Llama 3.1 70B Instruct", ContextLen: 131072, Description: "Mid-size Meta model", Category: []string{"general"}},
		{ID: "bedrock/meta.llama3-1-8b-instruct", Object: "model", OwnedBy: "meta", Provider: "bedrock", Name: "Llama 3.1 8B Instruct", ContextLen: 131072, Description: "Efficient Meta model", Category: []string{"general"}},
		{ID: "bedrock/mistral.mistral-large-2407", Object: "model", OwnedBy: "mistral", Provider: "bedrock", Name: "Mistral Large", ContextLen: 128000, Description: "Mistral's flagship model", Category: []string{"general"}, IsFlagship: true},
		{ID: "bedrock/mistral.mixtral-8x7b-instruct", Object: "model", OwnedBy: "mistral", Provider: "bedrock", Name: "Mixtral 8x7B", ContextLen: 32000, Description: "Efficient MoE model", Category: []string{"general"}},
		{ID: "bedrock/amazon.titan-text-express", Object: "model", OwnedBy: "amazon", Provider: "bedrock", Name: "Titan Text Express", ContextLen: 8192, Description: "Amazon's Titan model", Category: []string{"general"}},
		{ID: "bedrock/ai21.jamba-1-5-large", Object: "model", OwnedBy: "ai21", Provider: "bedrock", Name: "Jamba 1.5 Large", ContextLen: 256000, Description: "AI21's Jamba model", Category: []string{"general"}, IsFlagship: true},

		// OpenAI Models - Latest 2025-2026
		{ID: "openai/gpt-5.2", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-5.2", ContextLen: 256000, Description: "OpenAI's latest flagship model", Category: []string{"general", "coding", "reasoning"}, IsFlagship: true},
		{ID: "openai/gpt-4.5", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.5", ContextLen: 128000, Description: "OpenAI's previous flagship model", Category: []string{"general", "coding"}},
		{ID: "openai/gpt-4.1", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.1", ContextLen: 200000, Description: "High-capability instruction following", Category: []string{"general", "coding"}},
		{ID: "openai/gpt-4.1-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.1 Mini", ContextLen: 200000, Description: "Efficient instruction following", Category: []string{"general"}},
		{ID: "openai/gpt-4o", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o", ContextLen: 128000, Description: "Omni multimodal model", Category: []string{"general", "vision"}, IsFlagship: true},
		{ID: "openai/gpt-4o-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o Mini", ContextLen: 128000, Description: "Fast multimodal model", Category: []string{"general", "vision"}},
		{ID: "openai/o1", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1", ContextLen: 200000, Description: "Advanced reasoning model", Category: []string{"reasoning"}, IsFlagship: true},
		{ID: "openai/o1-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1 Mini", ContextLen: 128000, Description: "Fast reasoning model", Category: []string{"reasoning"}},
		{ID: "openai/o3", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o3", ContextLen: 200000, Description: "Next-gen reasoning model", Category: []string{"reasoning"}, IsFlagship: true},
		{ID: "openai/o3-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o3 Mini", ContextLen: 200000, Description: "Efficient reasoning model", Category: []string{"reasoning"}},
		{ID: "openai/o4-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o4 Mini", ContextLen: 200000, Description: "Compact reasoning model", Category: []string{"reasoning"}},
		{ID: "openai/gpt-4-turbo", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4 Turbo", ContextLen: 128000, Description: "Fast GPT-4 variant", Category: []string{"general", "coding"}},
		{ID: "openai/gpt-4", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4", ContextLen: 8192, Description: "Previous flagship model", Category: []string{"general", "coding"}},
		{ID: "openai/gpt-3.5-turbo", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-3.5 Turbo", ContextLen: 16385, Description: "Fast and efficient model", Category: []string{"general"}},
		{ID: "openai/gpt-oss-120b", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-OSS 120B", ContextLen: 128000, Description: "Open-weight flagship model", Category: []string{"general", "open_weights"}, IsFlagship: true},
		{ID: "openai/gpt-oss-20b", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-OSS 20B", ContextLen: 128000, Description: "Open-weight efficient model", Category: []string{"general", "open_weights"}},

		// xAI Grok Models - Latest 2025-2026
		{ID: "xai/grok-4", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4", ContextLen: 256000, Description: "xAI's flagship reasoning model", Category: []string{"general", "reasoning"}, IsFlagship: true},
		{ID: "xai/grok-4-fast", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4 Fast", ContextLen: 2000000, Description: "Fast Grok 4 variant", Category: []string{"general"}},
		{ID: "xai/grok-4.1", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4.1", ContextLen: 256000, Description: "Enhanced reasoning model", Category: []string{"general", "reasoning"}},
		{ID: "xai/grok-4.1-fast-reasoning", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4.1 Fast Reasoning", ContextLen: 2000000, Description: "Fast reasoning model", Category: []string{"reasoning"}},
		{ID: "xai/grok-3", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 3", ContextLen: 131072, Description: "Previous Grok flagship", Category: []string{"general"}},
		{ID: "xai/grok-3-mini", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 3 Mini", ContextLen: 131072, Description: "Efficient Grok model", Category: []string{"general"}},
		{ID: "xai/grok-2", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 2", ContextLen: 131072, Description: "Earlier Grok model", Category: []string{"general"}},
		{ID: "xai/grok-2-vision", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 2 Vision", ContextLen: 32768, Description: "Vision-enabled Grok", Category: []string{"vision"}},

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
		{ID: "anthropic/claude-opus-4-6", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Opus 4.6", ContextLen: 200000, Description: "Anthropic's flagship model", Category: []string{"general", "reasoning", "coding"}, IsFlagship: true},
		{ID: "anthropic/claude-opus-4-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Opus 4.5", ContextLen: 200000, Description: "High-capability model", Category: []string{"general", "reasoning"}},
		{ID: "anthropic/claude-opus-4-1", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Opus 4.1", ContextLen: 200000, Description: "Reasoning-focused model", Category: []string{"reasoning"}},
		{ID: "anthropic/claude-sonnet-4-6", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 4.6", ContextLen: 200000, Description: "Balanced performance model", Category: []string{"general", "coding"}, IsFlagship: true},
		{ID: "anthropic/claude-sonnet-4-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 4.5", ContextLen: 200000, Description: "Efficient balanced model", Category: []string{"general"}},
		{ID: "anthropic/claude-sonnet-4", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 4", ContextLen: 200000, Description: "Balanced performance", Category: []string{"general"}},
		{ID: "anthropic/claude-sonnet-3-7", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Sonnet 3.7", ContextLen: 200000, Description: "Previous balanced model", Category: []string{"general"}},
		{ID: "anthropic/claude-haiku-4-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Haiku 4.5", ContextLen: 200000, Description: "Fast efficient model", Category: []string{"general"}},
		{ID: "anthropic/claude-haiku-3-5", Object: "model", OwnedBy: "anthropic", Provider: "anthropic", Name: "Claude Haiku 3.5", ContextLen: 200000, Description: "Very fast model", Category: []string{"general"}},

		// Google Vertex AI Models - Latest 2025-2026
		{ID: "vertex/gemini-2.5-pro", Object: "model", OwnedBy: "google", Provider: "vertex", Name: "Gemini 2.5 Pro (Vertex)", ContextLen: 1000000, Description: "Google's flagship multimodal model on Vertex", Category: []string{"general", "vision", "multimodal"}, IsFlagship: true},
		{ID: "vertex/gemini-2.5-flash", Object: "model", OwnedBy: "google", Provider: "vertex", Name: "Gemini 2.5 Flash (Vertex)", ContextLen: 1000000, Description: "Fast multimodal model on Vertex", Category: []string{"general", "vision"}, IsFlagship: true},
		{ID: "vertex/gemini-2.0-pro", Object: "model", OwnedBy: "google", Provider: "vertex", Name: "Gemini 2.0 Pro (Vertex)", ContextLen: 2000000, Description: "High-capability Gemini 2.0 on Vertex", Category: []string{"general", "vision"}},
		{ID: "vertex/gemini-2.0-flash", Object: "model", OwnedBy: "google", Provider: "vertex", Name: "Gemini 2.0 Flash (Vertex)", ContextLen: 1000000, Description: "Fast Gemini 2.0 on Vertex", Category: []string{"general", "vision"}},
		{ID: "vertex/gemini-1.5-pro", Object: "model", OwnedBy: "google", Provider: "vertex", Name: "Gemini 1.5 Pro (Vertex)", ContextLen: 2000000, Description: "High-capability multimodal on Vertex", Category: []string{"general", "vision", "multimodal"}},

		// Google Gemini Models - Latest 2025-2026
		{ID: "gemini/gemini-3-pro-preview", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 3 Pro Preview", ContextLen: 1048576, Description: "Next-gen Gemini flagship", Category: []string{"general", "vision", "multimodal"}, IsFlagship: true},
		{ID: "gemini/gemini-3-flash-preview", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 3 Flash Preview", ContextLen: 1048576, Description: "Fast next-gen model", Category: []string{"general", "vision"}},
		{ID: "gemini/gemini-2.5-flash", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.5 Flash", ContextLen: 1000000, Description: "Fast multimodal model", Category: []string{"general", "vision", "multimodal"}, IsFlagship: true},
		{ID: "gemini/gemini-2.5-pro", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.5 Pro", ContextLen: 1000000, Description: "High-capability multimodal", Category: []string{"general", "vision", "multimodal"}, IsFlagship: true},
		{ID: "gemini/gemini-2.0-flash", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.0 Flash", ContextLen: 1000000, Description: "Fast Gemini 2.0", Category: []string{"general", "vision"}},
		{ID: "gemini/gemini-2.0-pro", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.0 Pro", ContextLen: 2000000, Description: "High-capability Gemini 2.0", Category: []string{"general", "vision"}},
		{ID: "gemini/gemini-2.0-flash-lite", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 2.0 Flash Lite", ContextLen: 1000000, Description: "Very efficient model", Category: []string{"general"}},
		{ID: "gemini/gemini-1.5-flash", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 1.5 Flash", ContextLen: 1000000, Description: "Fast multimodal", Category: []string{"general", "vision"}},
		{ID: "gemini/gemini-1.5-pro", Object: "model", OwnedBy: "google", Provider: "gemini", Name: "Gemini 1.5 Pro", ContextLen: 2000000, Description: "High-capability multimodal", Category: []string{"general", "vision", "multimodal"}},

		// DeepSeek Models - 2025
		{ID: "deepseek/deepseek-chat", Object: "model", OwnedBy: "deepseek", Provider: "deepseek", Name: "DeepSeek V3 Chat", ContextLen: 128000, Description: "DeepSeek's flagship chat model", Category: []string{"general", "coding"}, IsFlagship: true},
		{ID: "deepseek/deepseek-reasoner", Object: "model", OwnedBy: "deepseek", Provider: "deepseek", Name: "DeepSeek R1 Reasoner", ContextLen: 128000, Description: "Advanced reasoning model", Category: []string{"reasoning", "coding"}, IsFlagship: true},

		// Mistral AI Models - Latest 2025-2026
		{ID: "mistral/mistral-large-3", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Mistral Large 3", ContextLen: 256000, Description: "Mistral's flagship MoE model (41B/675B params)", Category: []string{"general", "multimodal"}, IsFlagship: true},
		{ID: "mistral/mistral-medium-3-1", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Mistral Medium 3.1", ContextLen: 128000, Description: "Frontier mid-tier multimodal", Category: []string{"general"}},
		{ID: "mistral/mistral-small-3-2", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Mistral Small 3.2", ContextLen: 32000, Description: "Fast general-purpose", Category: []string{"general"}},
		{ID: "mistral/devstral-2", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Devstral 2", ContextLen: 128000, Description: "Code agent model", Category: []string{"coding"}},
		{ID: "mistral/magistral-1-2", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Magistral 1.2", ContextLen: 128000, Description: "Reasoning model", Category: []string{"reasoning"}},
		{ID: "mistral/codestral-2508", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Codestral 25.08", ContextLen: 128000, Description: "Code generation", Category: []string{"coding"}},
		{ID: "mistral/ministral-14b", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Ministral 14B", ContextLen: 128000, Description: "Edge/embedded model", Category: []string{"general"}},
		{ID: "mistral/ministral-8b", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Ministral 8B", ContextLen: 128000, Description: "Edge/embedded model", Category: []string{"general"}},
		{ID: "mistral/ministral-3b", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Ministral 3B", ContextLen: 128000, Description: "Edge/embedded model", Category: []string{"general"}},
		{ID: "mistral/pixtral-large", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Pixtral Large", ContextLen: 128000, Description: "Vision model", Category: []string{"vision"}},
		{ID: "mistral/voxtral", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Voxtral", ContextLen: 128000, Description: "Audio model", Category: []string{"audio"}},
		{ID: "mistral/nemo", Object: "model", OwnedBy: "mistral", Provider: "mistral", Name: "Nemo", ContextLen: 128000, Description: "Base model", Category: []string{"general"}},

		// Groq Models - Latest 2025-2026
		{ID: "groq/llama-3.3-70b-versatile", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 3.3 70B", ContextLen: 131072, Description: "Fast inference Llama 3.3", Category: []string{"general"}, IsFlagship: true},
		{ID: "groq/llama-3.1-8b-instant", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 3.1 8B", ContextLen: 131072, Description: "Very fast inference", Category: []string{"general"}},
		{ID: "groq/llama-4-maverick", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 4 Maverick", ContextLen: 131072, Description: "Latest Llama 4 flagship", Category: []string{"general"}, IsFlagship: true},
		{ID: "groq/llama-4-scout", Object: "model", OwnedBy: "meta", Provider: "groq", Name: "Llama 4 Scout", ContextLen: 131072, Description: "Efficient Llama 4", Category: []string{"general"}},
		{ID: "groq/openai-gpt-oss-120b", Object: "model", OwnedBy: "openai", Provider: "groq", Name: "GPT-OSS 120B", ContextLen: 131072, Description: "OpenAI open-weight 120B", Category: []string{"general", "open_weights"}, IsFlagship: true},
		{ID: "groq/openai-gpt-oss-20b", Object: "model", OwnedBy: "openai", Provider: "groq", Name: "GPT-OSS 20B", ContextLen: 131072, Description: "OpenAI open-weight 20B", Category: []string{"general", "open_weights"}},
		{ID: "groq/qwen-2.5-32b", Object: "model", OwnedBy: "alibaba", Provider: "groq", Name: "Qwen 2.5 32B", ContextLen: 131072, Description: "Fast Qwen 2.5", Category: []string{"general"}},
		{ID: "groq/qwen-2.5-coder-32b", Object: "model", OwnedBy: "alibaba", Provider: "groq", Name: "Qwen 2.5 Coder 32B", ContextLen: 131072, Description: "Code-focused Qwen", Category: []string{"coding"}},
		{ID: "groq/qwen-qwq-32b", Object: "model", OwnedBy: "alibaba", Provider: "groq", Name: "Qwen QwQ 32B", ContextLen: 131072, Description: "Reasoning Qwen model", Category: []string{"reasoning"}},
		{ID: "groq/deepseek-r1-distill-qwen-32b", Object: "model", OwnedBy: "deepseek", Provider: "groq", Name: "DeepSeek R1 Distill Qwen 32B", ContextLen: 131072, Description: "Distilled reasoning model", Category: []string{"reasoning"}},
		{ID: "groq/gemma2-9b-it", Object: "model", OwnedBy: "google", Provider: "groq", Name: "Gemma 2 9B", ContextLen: 8192, Description: "Fast Gemma model", Category: []string{"general"}},
		{ID: "groq/mistral-saba-24b", Object: "model", OwnedBy: "mistral", Provider: "groq", Name: "Mistral Saba 24B", ContextLen: 32000, Description: "Specialized model", Category: []string{"general"}},

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
		{ID: "perplexity/sonar", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar", ContextLen: 128000, Description: "Lightweight search-augmented model", Category: []string{"search", "general"}},
		{ID: "perplexity/sonar-pro", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar Pro", ContextLen: 200000, Description: "Advanced search-augmented model", Category: []string{"search", "general"}, IsFlagship: true},
		{ID: "perplexity/sonar-reasoning", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar Reasoning", ContextLen: 128000, Description: "Reasoning with search", Category: []string{"search", "reasoning"}},
		{ID: "perplexity/sonar-reasoning-pro", Object: "model", OwnedBy: "perplexity", Provider: "perplexity", Name: "Sonar Reasoning Pro", ContextLen: 128000, Description: "Advanced reasoning with search", Category: []string{"search", "reasoning"}, IsFlagship: true},

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

		// Ollama Models (Local Inference)
		{ID: "ollama/llama3.3", Object: "model", OwnedBy: "meta", Provider: "ollama", Name: "Llama 3.3", ContextLen: 131072},
		{ID: "ollama/llama3.2", Object: "model", OwnedBy: "meta", Provider: "ollama", Name: "Llama 3.2", ContextLen: 131072},
		{ID: "ollama/llama3.1", Object: "model", OwnedBy: "meta", Provider: "ollama", Name: "Llama 3.1", ContextLen: 131072},
		{ID: "ollama/mistral", Object: "model", OwnedBy: "mistral", Provider: "ollama", Name: "Mistral", ContextLen: 32768},
		{ID: "ollama/codellama", Object: "model", OwnedBy: "meta", Provider: "ollama", Name: "Code Llama", ContextLen: 16384},
		{ID: "ollama/qwen2.5", Object: "model", OwnedBy: "qwen", Provider: "ollama", Name: "Qwen 2.5", ContextLen: 131072},
		{ID: "ollama/deepseek-coder-v2", Object: "model", OwnedBy: "deepseek", Provider: "ollama", Name: "DeepSeek Coder V2", ContextLen: 128000},
		{ID: "ollama/phi4", Object: "model", OwnedBy: "microsoft", Provider: "ollama", Name: "Phi-4", ContextLen: 16384},
		{ID: "ollama/gemma2", Object: "model", OwnedBy: "google", Provider: "ollama", Name: "Gemma 2", ContextLen: 8192},
		{ID: "ollama/mixtral", Object: "model", OwnedBy: "mistral", Provider: "ollama", Name: "Mixtral 8x7B", ContextLen: 32768},

		// DeepInfra Models - Latest 2025-2026
		{ID: "deepinfra/llama-3.3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "deepinfra", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "deepinfra/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "deepinfra", Name: "DeepSeek V3", ContextLen: 163840},
		{ID: "deepinfra/qwen3-max-thinking", Object: "model", OwnedBy: "qwen", Provider: "deepinfra", Name: "Qwen3 Max Thinking", ContextLen: 131072},
		{ID: "deepinfra/kimi-k2.5", Object: "model", OwnedBy: "moonshot", Provider: "deepinfra", Name: "Kimi K2.5", ContextLen: 262144},
		{ID: "deepinfra/glm-5", Object: "model", OwnedBy: "z-ai", Provider: "deepinfra", Name: "GLM 5", ContextLen: 131072},
		{ID: "deepinfra/mistral-large", Object: "model", OwnedBy: "mistral", Provider: "deepinfra", Name: "Mistral Large", ContextLen: 131072},

		// SambaNova Models - Latest 2025-2026
		{ID: "sambanova/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "sambanova", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "sambanova/llama-3.2-90b-vision", Object: "model", OwnedBy: "meta", Provider: "sambanova", Name: "Llama 3.2 90B Vision", ContextLen: 131072},
		{ID: "sambanova/qwen2.5-72b", Object: "model", OwnedBy: "qwen", Provider: "sambanova", Name: "Qwen 2.5 72B", ContextLen: 131072},

		// Nebius AI Models - Latest 2025-2026
		{ID: "nebius/llama-3.3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "nebius", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "nebius/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "nebius", Name: "DeepSeek V3", ContextLen: 128000},
		{ID: "nebius/qwen2.5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "nebius", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},

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

		// Moonshot AI (Kimi) Models - Latest 2025-2026
		{ID: "moonshot/kimi-k2.5", Object: "model", OwnedBy: "moonshot", Provider: "moonshot", Name: "Kimi K2.5", ContextLen: 262144},
		{ID: "moonshot/kimi-k2", Object: "model", OwnedBy: "moonshot", Provider: "moonshot", Name: "Kimi K2", ContextLen: 262144},
		{ID: "moonshot/kimi-k1.5", Object: "model", OwnedBy: "moonshot", Provider: "moonshot", Name: "Kimi K1.5", ContextLen: 262144},
		{ID: "moonshot/kimi-v1-8k", Object: "model", OwnedBy: "moonshot", Provider: "moonshot", Name: "Kimi v1 8K", ContextLen: 8192},
		{ID: "moonshot/kimi-v1-32k", Object: "model", OwnedBy: "moonshot", Provider: "moonshot", Name: "Kimi v1 32K", ContextLen: 32768},
		{ID: "moonshot/kimi-v1-128k", Object: "model", OwnedBy: "moonshot", Provider: "moonshot", Name: "Kimi v1 128K", ContextLen: 131072},

		// MiniMax Models - Latest 2025-2026
		{ID: "minimax/abab6.5", Object: "model", OwnedBy: "minimax", Provider: "minimax", Name: "abab 6.5", ContextLen: 200000},
		{ID: "minimax/abab6.5s", Object: "model", OwnedBy: "minimax", Provider: "minimax", Name: "abab 6.5s", ContextLen: 200000},
		{ID: "minimax/m2.5", Object: "model", OwnedBy: "minimax", Provider: "minimax", Name: "MiniMax M2.5", ContextLen: 204800},
		{ID: "minimax/m2.1", Object: "model", OwnedBy: "minimax", Provider: "minimax", Name: "MiniMax M2.1", ContextLen: 204800},

		// SiliconFlow Models - Latest 2025-2026
		{ID: "siliconflow/qwen-2.5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "siliconflow", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},
		{ID: "siliconflow/qwen-2.5-vl-7b-instruct", Object: "model", OwnedBy: "qwen", Provider: "siliconflow", Name: "Qwen 2.5 VL 7B Instruct", ContextLen: 32768},
		{ID: "siliconflow/qwen-3-32b", Object: "model", OwnedBy: "qwen", Provider: "siliconflow", Name: "Qwen 3 32B", ContextLen: 131072},
		{ID: "siliconflow/qwen-3-14b", Object: "model", OwnedBy: "qwen", Provider: "siliconflow", Name: "Qwen 3 14B", ContextLen: 131072},
		{ID: "siliconflow/deepseek-v2.5", Object: "model", OwnedBy: "deepseek", Provider: "siliconflow", Name: "DeepSeek V2.5", ContextLen: 128000},
		{ID: "siliconflow/deepseek-r1-distill-qwen-14b", Object: "model", OwnedBy: "deepseek", Provider: "siliconflow", Name: "DeepSeek R1 Distill Qwen 14B", ContextLen: 131072},
		{ID: "siliconflow/llama-3.1-8b-instruct", Object: "model", OwnedBy: "meta", Provider: "siliconflow", Name: "Llama 3.1 8B Instruct", ContextLen: 131072},
		{ID: "siliconflow/glm-4-9b-chat", Object: "model", OwnedBy: "zhipu", Provider: "siliconflow", Name: "GLM 4 9B Chat", ContextLen: 131072},

		// Replicate Models - Latest 2025-2026
		{ID: "replicate/llama-3.3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "replicate", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "replicate/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "replicate", Name: "DeepSeek V3", ContextLen: 128000},
		{ID: "replicate/qwen3-32b", Object: "model", OwnedBy: "qwen", Provider: "replicate", Name: "Qwen 3 32B", ContextLen: 131072},
		{ID: "replicate/gpt-5", Object: "model", OwnedBy: "openai", Provider: "replicate", Name: "GPT-5", ContextLen: 256000},
		{ID: "replicate/gemini-3-pro", Object: "model", OwnedBy: "google", Provider: "replicate", Name: "Gemini 3 Pro", ContextLen: 1048576},
		{ID: "replicate/claude-4.5-sonnet", Object: "model", OwnedBy: "anthropic", Provider: "replicate", Name: "Claude 4.5 Sonnet", ContextLen: 200000},

		// Anyscale Models - Latest 2025-2026
		{ID: "anyscale/llama-3.3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "anyscale", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "anyscale/llama-3.1-8b-instruct", Object: "model", OwnedBy: "meta", Provider: "anyscale", Name: "Llama 3.1 8B Instruct", ContextLen: 131072},
		{ID: "anyscale/llama-4-maverick", Object: "model", OwnedBy: "meta", Provider: "anyscale", Name: "Llama 4 Maverick", ContextLen: 131072},
		{ID: "anyscale/mistral-7b-instruct", Object: "model", OwnedBy: "mistral", Provider: "anyscale", Name: "Mistral 7B Instruct", ContextLen: 32768},

		// Cerebras Models - Latest 2025-2026
		{ID: "cerebras/llama-4-scout", Object: "model", OwnedBy: "meta", Provider: "cerebras", Name: "Llama 4 Scout", ContextLen: 131072},
		{ID: "cerebras/llama-4-maverick", Object: "model", OwnedBy: "meta", Provider: "cerebras", Name: "Llama 4 Maverick", ContextLen: 131072},
		{ID: "cerebras/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "cerebras", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "cerebras/llama-3.1-405b", Object: "model", OwnedBy: "meta", Provider: "cerebras", Name: "Llama 3.1 405B", ContextLen: 131072},

		// Baseten Models - Latest 2025-2026
		{ID: "baseten/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "baseten", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "baseten/qwen-2.5-72b", Object: "model", OwnedBy: "qwen", Provider: "baseten", Name: "Qwen 2.5 72B", ContextLen: 131072},

		// AI21 Labs Models - Latest 2025-2026
		{ID: "ai21/jamba-2", Object: "model", OwnedBy: "ai21", Provider: "ai21", Name: "Jamba 2", ContextLen: 256000},
		{ID: "ai21/jamba-1.5", Object: "model", OwnedBy: "ai21", Provider: "ai21", Name: "Jamba 1.5", ContextLen: 256000},

		// NVIDIA NIM Models - Latest 2025-2026
		{ID: "nvidia/llama-3.3-70b-instruct", Object: "model", OwnedBy: "meta", Provider: "nvidia", Name: "Llama 3.3 70B Instruct", ContextLen: 131072},
		{ID: "nvidia/llama-3.1-405b-instruct", Object: "model", OwnedBy: "meta", Provider: "nvidia", Name: "Llama 3.1 405B Instruct", ContextLen: 131072},
		{ID: "nvidia/qwen2.5-72b-instruct", Object: "model", OwnedBy: "qwen", Provider: "nvidia", Name: "Qwen 2.5 72B Instruct", ContextLen: 131072},

		// FriendliAI Models - Latest 2025-2026
		{ID: "friendli/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "friendli", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "friendli/deepseek-v3", Object: "model", OwnedBy: "deepseek", Provider: "friendli", Name: "DeepSeek V3", ContextLen: 128000},
		{ID: "friendli/qwen-2.5-72b", Object: "model", OwnedBy: "qwen", Provider: "friendli", Name: "Qwen 2.5 72B", ContextLen: 131072},

		// Venice AI Models - Latest 2025-2026
		{ID: "venice/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "venice", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "venice/qwen-2.5-72b", Object: "model", OwnedBy: "qwen", Provider: "venice", Name: "Qwen 2.5 72B", ContextLen: 131072},

		// OVHCloud AI Models - Latest 2025-2026
		{ID: "ovhcloud/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "ovhcloud", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "ovhcloud/qwen-2.5-72b", Object: "model", OwnedBy: "qwen", Provider: "ovhcloud", Name: "Qwen 2.5 72B", ContextLen: 131072},

		// Scaleway AI Models - Latest 2025-2026
		{ID: "scaleway/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "scaleway", Name: "Llama 3.3 70B", ContextLen: 131072},
		{ID: "scaleway/qwen-2.5-72b", Object: "model", OwnedBy: "qwen", Provider: "scaleway", Name: "Qwen 2.5 72B", ContextLen: 131072},

		// StepFun Models - Latest 2025-2026
		{ID: "stepfun/step-1v", Object: "model", OwnedBy: "stepfun", Provider: "stepfun", Name: "Step 1V", ContextLen: 131072},
		{ID: "stepfun/step-1", Object: "model", OwnedBy: "stepfun", Provider: "stepfun", Name: "Step 1", ContextLen: 131072},

		// Xiaomi MiMo Models - Latest 2025-2026
		{ID: "xiaomi/mimo", Object: "model", OwnedBy: "xiaomi", Provider: "xiaomi", Name: "MiMo", ContextLen: 131072},

		// Liquid AI Models - Latest 2025-2026
		{ID: "liquid/liquid-1", Object: "model", OwnedBy: "liquid", Provider: "liquid", Name: "Liquid 1", ContextLen: 131072},

		// Arcee AI Models
		{ID: "arcee/arity-falcon", Object: "model", OwnedBy: "arcee", Provider: "arcee", Name: "Arity Falcon", ContextLen: 131072, Description: "Enterprise-focused language model", Category: []string{"general", "enterprise"}, IsFlagship: true},

		// Chutes AI Models
		{ID: "chutes/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "chutes", Name: "Llama 3.3 70B", ContextLen: 131072, Description: "Multilingual large language model", Category: []string{"general"}, IsFlagship: true},

		// Morph AI Models
		{ID: "morph/llama-3.1-70b", Object: "model", OwnedBy: "meta", Provider: "morph", Name: "Llama 3.1 70B", ContextLen: 131072, Description: "Instruction-following conversational model", Category: []string{"general"}},

		// NextBit Models
		{ID: "nextbit/llama-3.1-8b", Object: "model", OwnedBy: "meta", Provider: "nextbit", Name: "Llama 3.1 8B", ContextLen: 131072, Description: "Efficient instruction-following model", Category: []string{"general"}},

		// Parasail Models
		{ID: "parasail/llama-3.3-70b", Object: "model", OwnedBy: "meta", Provider: "parasail", Name: "Llama 3.3 70B", ContextLen: 131072, Description: "High-performance multilingual model", Category: []string{"general"}, IsFlagship: true},

		// Phala Network Models
		{ID: "phala/llama-3.1-70b", Object: "model", OwnedBy: "meta", Provider: "phala", Name: "Llama 3.1 70B", ContextLen: 131072, Description: "Privacy-focused inference", Category: []string{"general"}},

		// ModelRun Models
		{ID: "modelrun/llama-3.1-8b", Object: "model", OwnedBy: "meta", Provider: "modelrun", Name: "Llama 3.1 8B", ContextLen: 131072, Description: "Fast inference model", Category: []string{"general"}},

		// Cloudflare Workers AI Models - Latest 2025-2026
		{ID: "cloudflare/llama-3.1-8b", Object: "model", OwnedBy: "meta", Provider: "cloudflare", Name: "Llama 3.1 8B", ContextLen: 128000},
		{ID: "cloudflare/llama-3.2-1b", Object: "model", OwnedBy: "meta", Provider: "cloudflare", Name: "Llama 3.2 1B", ContextLen: 128000},
		{ID: "cloudflare/qwen-2.5-7b", Object: "model", OwnedBy: "qwen", Provider: "cloudflare", Name: "Qwen 2.5 7B", ContextLen: 128000},
		{ID: "cloudflare/gemma-2-2b", Object: "model", OwnedBy: "google", Provider: "cloudflare", Name: "Gemma 2 2B", ContextLen: 128000},
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
