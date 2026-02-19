package vertex

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

type Provider struct {
	projectID  string
	location   string
	httpClient *http.Client
}

func New(projectID, location string) *Provider {
	if projectID == "" {
		projectID = os.Getenv("GOOGLE_PROJECT_ID")
	}
	if location == "" {
		location = os.Getenv("GOOGLE_LOCATION")
		if location == "" {
			location = "us-central1"
		}
	}
	return &Provider{
		projectID: projectID,
		location:  location,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *Provider) Name() string { return "vertex" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	reqBody := convertToVertexRequest(req)
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	modelName := req.Model
	if strings.HasPrefix(modelName, "vertex/") {
		modelName = strings.TrimPrefix(modelName, "vertex/")
	}
	if strings.HasPrefix(modelName, "google/") {
		modelName = strings.TrimPrefix(modelName, "google/")
	}

	url := fmt.Sprintf("https://%s-aiplatform.googleapis.com/v1/projects/%s/locations/%s/endpoints/%s/openapi",
		p.location, p.projectID, p.location, modelName)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	token := os.Getenv("GOOGLE_ACCESS_TOKEN")
	if token == "" {
		token = os.Getenv("GOOGLE_OAUTH_TOKEN")
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("vertex API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var vertexResp VertexResponse
	if err := json.Unmarshal(respBody, &vertexResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return convertFromVertexResponse(vertexResp), nil
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	return fmt.Errorf("streaming not yet implemented for Vertex AI")
}
