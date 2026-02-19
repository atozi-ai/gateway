package aws_bedrock

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/atozi-ai/gateway/internal/domain/llm"
)

const (
	awsService = "bedrock-runtime"
	awsRegion  = "us-east-1"
)

type Provider struct {
	awsAccessKey string
	awsSecretKey string
	awsRegion    string
	httpClient   *http.Client
}

func New(accessKey, secretKey, region string) *Provider {
	if region == "" {
		region = awsRegion
	}
	return &Provider{
		awsAccessKey: accessKey,
		awsSecretKey: secretKey,
		awsRegion:    region,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}
}

func (p *Provider) Name() string { return "bedrock" }

func (p *Provider) Chat(ctx context.Context, req llm.ChatRequest) (*llm.ChatResponse, error) {
	bedrockReq := convertToBedrockRequest(req)
	bedrockReq.Model = req.Model

	body, err := json.Marshal(bedrockReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com/model/%s/converse", p.awsRegion, req.Model)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	p.signRequest(httpReq, body)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

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
		return nil, fmt.Errorf("bedrock API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var bedrockResp ConverseResponse
	if err := json.Unmarshal(respBody, &bedrockResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return convertFromBedrockResponse(bedrockResp), nil
}

func (p *Provider) ChatStream(ctx context.Context, req llm.ChatRequest, callback func(*llm.StreamChunk) error) error {
	bedrockReq := convertToBedrockRequest(req)
	bedrockReq.Model = req.Model

	body, err := json.Marshal(bedrockReq)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("https://bedrock-runtime.%s.amazonaws.com/model/%s/converse-stream", p.awsRegion, req.Model)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	p.signRequest(httpReq, body)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk ConverseStreamChunk
		if err := decoder.Decode(&chunk); err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to decode stream: %w", err)
		}

		if chunk.Delta != nil && chunk.Delta.Text != "" {
			content := chunk.Delta.Text
			streamChunk := &llm.StreamChunk{
				Choices: []llm.StreamChoice{
					{
						Delta: llm.StreamDelta{
							Content: &content,
						},
					},
				},
			}
			if err := callback(streamChunk); err != nil {
				return err
			}
		}

		if chunk.StopReason != "" {
			break
		}
	}

	return nil
}

func (p *Provider) signRequest(req *http.Request, body []byte) {
	now := time.Now().UTC()
	date := now.Format("20060102T150405Z")
	shortDate := now.Format("20060102")

	req.Header.Set("X-Amz-Date", date)
	req.Header.Set("X-Amz-Target", "AmazonBedrockRuntime.Converse")

	host := req.URL.Host
	contentHash := sha256Hash(body)
	req.Header.Set("X-Amz-Content-Sha256", contentHash)

	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-amz-content-sha256:%s\nx-amz-date:%s\nx-amz-target:%s\n",
		req.Header.Get("Content-Type"), host, contentHash, date, req.Header.Get("X-Amz-Target"))

	signedHeaders := "content-type;host;x-amz-content-sha256;x-amz-date;x-amz-target"
	canonicalRequest := fmt.Sprintf("POST\n%s\n/model/%s/converse\n\n%s\n\n%s\n%s",
		canonicalHeaders, strings.Split(req.URL.Path, "/")[3], signedHeaders, contentHash, "aws4_request")

	hashedCanonicalRequest := sha256Hash([]byte(canonicalRequest))
	credentialScope := fmt.Sprintf("%s/%s/%s/aws4_request", shortDate, p.awsRegion, awsService)
	stringToSign := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s\n%s", date, credentialScope, hashedCanonicalRequest)

	signingKey := getSignatureKey(p.awsSecretKey, shortDate, p.awsRegion, awsService)
	signature := hmacSHA256(signingKey, []byte(stringToSign))

	authHeader := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		p.awsAccessKey, credentialScope, signedHeaders, hex.EncodeToString(signature))
	req.Header.Set("Authorization", authHeader)
}

func sha256Hash(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func getSignatureKey(secretKey, dateStamp, regionName, serviceName string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secretKey), []byte(dateStamp))
	kRegion := hmacSHA256(kDate, []byte(regionName))
	kService := hmacSHA256(kRegion, []byte(serviceName))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	return kSigning
}
