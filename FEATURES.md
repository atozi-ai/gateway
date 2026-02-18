# Features

## Overview

AI Gateway for routing requests across multiple LLM providers with built-in resilience patterns.

## Core Features

### Multi-Provider Support
- **OpenAI** - `openai/<model>` (e.g., `openai/gpt-4o`)
- **Azure OpenAI** - `azure/<deployment>` (requires endpoint)
- **xAI (Grok)** - `xai/<model>` (e.g., `xai/grok-2`)
- **Z.ai** - `zai/<model>`

### Chat Completions

#### Request Options
- `temperature` - Sampling temperature
- `max_tokens` - Maximum tokens to generate
- `top_p` - Nucleus sampling
- `stop` - Stop sequences
- `frequency_penalty` - Frequency penalty
- `presence_penalty` - Presence penalty
- `logit_bias` - Logit bias for tokens
- `logprobs` - Return log probabilities
- `top_logprobs` - Number of top logprobs
- `n` - Number of completions
- `seed` - Random seed
- `user` - User identifier
- `verbosity` - Response verbosity (low/medium/high)

#### Structured Output
- `response_format` - JSON schema validation
- `response_format.type` - "json_schema" or "json_object"
- `response_format.schema` - JSON schema definition

#### Tool Calling
- `tools` - Function definitions for tool calling
- `tool_choice` - Force specific tool or auto
- `parallel_tool_calls` - Enable parallel tool execution
- `tool_resolution` - Tool resolution strategy

#### Streaming
- `stream` - Enable Server-Sent Events streaming
- `stream_options.include_usage` - Include token usage in chunks
- `stream_options.include_accumulated` - Include accumulated content in chunks

#### Response Options
- `raw` - Include raw provider response
- `include_accumulated` - Include accumulated content field

### Request Validation
- Required: `model`, `messages`
- Max messages: 1000
- Max request body: 10MB
- Request timeout: 120 seconds (configurable)

## Resilience Features

### Rate Limiting
- Per-API key rate limiting
- Multi-window support:
  - Per-second (token bucket with burst)
  - Per-minute
  - Per-hour
  - Per-day
- Configurable via environment variables
- Automatic cleanup of stale entries

### Circuit Breaker
- Automatic failure detection
- Opens after 5 failures with 50%+ failure ratio
- Closes after 3 successes
- 30-second timeout before retry
- Logs warnings on circuit state changes
- Per-provider circuit isolation

### Connection Pooling
- Shared HTTP client across all providers
- Max idle connections: 200
- Max idle per host: 50
- Connection timeout: 90 seconds
- TLS handshake timeout: 10 seconds

### Provider Caching
- Provider instances cached by (provider + apiKey + endpoint)
- Reduces memory allocation for repeated requests

## API

### Endpoints

#### Health Check
```
GET /health
```
Returns: `OK`

#### Chat Completions
```
POST /api/v1/chat/completions
Authorization: Bearer <api_key>
```

### Error Handling
- Standardized error responses
- Provider-specific error forwarding
- HTTP status codes: 400, 401, 429, 500, 502, 503

## Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | 8082 | Server port |
| `RATE_LIMIT_REQUESTS_PER_SECOND` | 10 | Requests per second per API key |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | 0 | Requests per minute (0=disabled) |
| `RATE_LIMIT_REQUESTS_PER_HOUR` | 0 | Requests per hour (0=disabled) |
| `RATE_LIMIT_REQUESTS_PER_DAY` | 0 | Requests per day (0=disabled) |
| `RATE_LIMIT_BURST` | 20 | Burst size for token bucket |
| `CIRCUIT_BREAKER_FAILURE_THRESHOLD` | 5 | Failures before circuit opens |
| `CIRCUIT_BREAKER_SUCCESS_THRESHOLD` | 3 | Successes to close circuit |
| `CIRCUIT_BREAKER_TIMEOUT_SECONDS` | 30 | Circuit open duration |

## Architecture

### Layers
1. **HTTP Handler** - Request parsing, validation
2. **Chat Handler** - Business logic, response formatting
3. **Provider Registry** - Provider selection, caching
4. **Circuit Breaker** - Failure isolation
5. **HTTP Client** - Connection pooling
6. **LLM Providers** - OpenAI, Azure, xAI, ZAI

### Request Flow
```
HTTP Request → Rate Limit → Handler → Provider Registry
           → Circuit Breaker → HTTP Client → Provider
           → Response → Client
```
