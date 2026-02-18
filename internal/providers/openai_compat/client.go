package openaicompat

import (
	"net/http"
	"sync"
	"time"
)

var (
	sharedHTTPClient *http.Client
	once             sync.Once
)

func initSharedClient() {
	transport := &http.Transport{
		MaxIdleConns:        200,
		MaxIdleConnsPerHost: 50,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
	}
	sharedHTTPClient = &http.Client{
		Timeout:   120 * time.Second,
		Transport: transport,
	}
}

func getSharedClient() *http.Client {
	once.Do(initSharedClient)
	return sharedHTTPClient
}

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
// Uses a shared HTTP client with connection pooling for better performance.
func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: getSharedClient(),
	}
}

// NewClientWithCustomHTTP returns a Client with a custom HTTP client.
// Use this when you need different timeout or transport settings.
func NewClientWithCustomHTTP(cfg Config, httpClient *http.Client) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

// setHeaders applies standard and custom headers to the request.
// If the request has an APIKey, it will be used instead of the client's default.
func (c *Client) setHeaders(req *http.Request, apiKey string) {
	// Use API key from request if provided, otherwise use client's default
	key := apiKey
	if key == "" {
		key = c.cfg.APIKey
	}

	// Skip Authorization header if api-key header is provided (e.g., Azure)
	if _, hasAPIKey := c.cfg.Headers["api-key"]; !hasAPIKey {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	req.Header.Set("Content-Type", "application/json")
	for k, v := range c.cfg.Headers {
		req.Header.Set(k, v)
	}
}
