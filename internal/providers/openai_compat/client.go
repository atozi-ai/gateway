package openaicompat

import (
	"net/http"
	"time"
)

// Config holds settings for an OpenAI-compatible API client.
type Config struct {
	BaseURL string
	APIKey  string
	Headers map[string]string // Extra headers beyond Authorization and Content-Type.
}

// Client performs HTTP calls against an OpenAI-compatible chat completions API.
type Client struct {
	cfg        Config
	httpClient *http.Client
}

// NewClient returns a ready-to-use Client for the given configuration.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: newHTTPClient(),
	}
}

func newHTTPClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	return &http.Client{
		Timeout:   60 * time.Second,
		Transport: transport,
	}
}

// setHeaders applies standard and custom headers to the request.
func (c *Client) setHeaders(req *http.Request) {
	// Skip Authorization header if api-key header is provided (e.g., Azure)
	if _, hasAPIKey := c.cfg.Headers["api-key"]; !hasAPIKey {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.cfg.Headers {
		req.Header.Set(k, v)
	}
}
