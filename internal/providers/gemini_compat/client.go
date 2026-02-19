package gemini_compat

import (
	"net/http"
	"sync"
	"time"
)

const (
	baseURL = "https://generativelanguage.googleapis.com"
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

type Config struct {
	BaseURL string
	APIKey  string
}

type Client struct {
	cfg        Config
	httpClient *http.Client
}

func NewClient(cfg Config) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: getSharedClient(),
	}
}

func NewClientWithCustomHTTP(cfg Config, httpClient *http.Client) *Client {
	return &Client{
		cfg:        cfg,
		httpClient: httpClient,
	}
}

func (c *Client) setHeaders(req *http.Request, apiKey string) {
	key := apiKey
	if key == "" {
		key = c.cfg.APIKey
	}

	req.Header.Set("x-goog-api-key", key)
	req.Header.Set("Content-Type", "application/json")
}
