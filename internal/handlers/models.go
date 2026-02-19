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
