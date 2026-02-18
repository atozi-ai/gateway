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
		{ID: "openai/gpt-4.1", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.1", ContextLen: 200000},
		{ID: "openai/gpt-4.1-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.1 Mini", ContextLen: 200000},
		{ID: "openai/gpt-4o", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o", ContextLen: 128000},
		{ID: "openai/gpt-4o-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o Mini", ContextLen: 128000},
		{ID: "openai/gpt-4o-realtime", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4o Realtime", ContextLen: 128000},
		{ID: "openai/o1", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1", ContextLen: 200000},
		{ID: "openai/o1-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1 Mini", ContextLen: 100000},
		{ID: "openai/o1-preview", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o1 Preview", ContextLen: 128000},
		{ID: "openai/o3", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o3", ContextLen: 200000},
		{ID: "openai/o3-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o3 Mini", ContextLen: 200000},
		{ID: "openai/o4-mini", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "o4 Mini", ContextLen: 200000},
		{ID: "openai/gpt-4-turbo", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4 Turbo", ContextLen: 128000},
		{ID: "openai/gpt-4", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4", ContextLen: 128000},
		{ID: "openai/gpt-3.5-turbo", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-3.5 Turbo", ContextLen: 16385},
		{ID: "openai/gpt-4.5", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "GPT-4.5", ContextLen: 128000},
		{ID: "openai/text-embedding-3-small", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "Embedding 3 Small", ContextLen: 8192},
		{ID: "openai/text-embedding-3-large", Object: "model", OwnedBy: "openai", Provider: "openai", Name: "Embedding 3 Large", ContextLen: 8192},
		{ID: "xai/grok-4", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4", ContextLen: 256000},
		{ID: "xai/grok-4-fast-reasoning", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 4 Fast Reasoning", ContextLen: 2000000},
		{ID: "xai/grok-3", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 3", ContextLen: 131072},
		{ID: "xai/grok-3-mini", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 3 Mini", ContextLen: 131072},
		{ID: "xai/grok-2", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 2", ContextLen: 131072},
		{ID: "xai/grok-2-vision-1212", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok 2 Vision", ContextLen: 32768},
		{ID: "xai/grok-beta", Object: "model", OwnedBy: "xai", Provider: "xai", Name: "Grok Beta", ContextLen: 131072},
		{ID: "zai/zai-core", Object: "model", OwnedBy: "zai", Provider: "zai", Name: "ZAI Core", ContextLen: 128000},
		{ID: "zai/zai-core-2025-01-20", Object: "model", OwnedBy: "zai", Provider: "zai", Name: "ZAI Core 2025-01-20", ContextLen: 128000},
		{ID: "zai/zai-core-2024-12-12", Object: "model", OwnedBy: "zai", Provider: "zai", Name: "ZAI Core 2024-12-12", ContextLen: 128000},
		{ID: "azure/gpt-4o", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4o", ContextLen: 128000},
		{ID: "azure/gpt-4o-mini", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4o Mini", ContextLen: 128000},
		{ID: "azure/gpt-4-turbo", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4 Turbo", ContextLen: 128000},
		{ID: "azure/gpt-4", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-4", ContextLen: 8192},
		{ID: "azure/gpt-35-turbo", Object: "model", OwnedBy: "azure", Provider: "azure", Name: "Azure GPT-3.5 Turbo", ContextLen: 16385},
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
