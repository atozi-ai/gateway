# Atozi AI Gateway

A unified AI Gateway - a single entry point for ALL AI services. Instead of managing multiple API keys, dealing with different provider APIs, and building observability from scratch, connect to ONE gateway and get access to everything.

**Currently Supported:**
- LLM/Chat - OpenAI, Anthropic, Google Gemini, AWS Bedrock, Azure OpenAI, Vertex AI, and 40+ more providers

**Coming Soon:**
- Embeddings, Speech, Vision, Document Extraction, Search, Moderation, RAG, and more

## Features

- **46+ LLM Providers** - OpenAI, Anthropic, Google Gemini, AWS Bedrock, Azure OpenAI, Vertex AI, and many more
- **Unified API** - OpenAI-compatible chat completions interface
- **Provider Options** - Pass provider-specific credentials (AWS keys, GCP project, Azure endpoint) via request options
- **Streaming Support** - Full streaming support across all providers
- **Structured Output** - JSON schema validation for typed responses
- **Tool Calling** - Function calling capability
- **Rate Limiting** - Configurable rate limits per second/minute/hour/day
- **Circuit Breaker** - Automatic failover on provider failures
- **Retry with Fallback** - Automatic retries with fallback to alternative models

## Installation

### Step 1: Download the Latest Release

**Linux/macOS:**
```bash
# Download the latest release (replace with actual version)
curl -L -o atozi-gateway https://github.com/atozi-ai/gateway/releases/latest/download/atozi-gateway-linux-amd64

# Make it executable
chmod +x atozi-gateway
```

**macOS (Apple Silicon):**
```bash
curl -L -o atozi-gateway https://github.com/atozi-ai/gateway/releases/latest/download/atozi-gateway-darwin-arm64
chmod +x atozi-gateway
```

**Windows:**
```powershell
# Using PowerShell
Invoke-WebRequest -Uri "https://github.com/atozi-ai/gateway/releases/latest/download/atozi-gateway-windows-amd64.exe" -OutFile "atozi-gateway.exe"
```

### Step 2: Setup Environment Variables

Copy the example environment file and configure as needed:

```bash
cp .env.example .env
```

Edit `.env` with your preferred settings:

```bash
# Rate Limiting (optional, defaults shown)
RATE_LIMIT_REQUESTS_PER_SECOND=10
RATE_LIMIT_REQUESTS_PER_MINUTE=500
RATE_LIMIT_REQUESTS_PER_HOUR=10000
RATE_LIMIT_REQUESTS_PER_DAY=100000
RATE_LIMIT_BURST=20
RATE_LIMIT_MAX_CLIENTS=10000

# Server Port (optional, default: 8082)
PORT=8082
```

### Step 3: Run the Gateway

```bash
# Linux/macOS
./atozi-gateway

# Windows
atozi-gateway.exe
```

The API will be available at `http://localhost:8082`

---

## API Endpoints

### Get Models List

**Endpoint:** `GET /api/v1/models`

**Request:**
```bash
curl -X GET http://localhost:8082/api/v1/models \
  -H "Authorization: Bearer YOUR_API_KEY"
```

**Response:**
```json
{
  "object": "list",
  "data": [
    {
      "id": "gpt-4o",
      "object": "model",
      "created": 1715367049,
      "owned_by": "openai",
      "provider": "openai",
      "description": "Latest GPT-4o model",
      "category": ["reasoning", "coding"],
      "is_flagship": true
    },
    {
      "id": "claude-3-5-sonnet-20241022",
      "object": "model",
      "created": 1715367049,
      "owned_by": "anthropic",
      "provider": "anthropic",
      "description": "Claude 3.5 Sonnet",
      "category": ["reasoning", "coding"],
      "is_flagship": true
    }
  ]
}
```

---

### Chat Completions

All providers follow the OpenAI-compatible format. Use the model prefix to specify the provider.

#### 1. OpenAI and OpenAI-Compatible Providers

**Supported:** openai, deepseek, mistral, groq, together, fireworks, perplexity, cohere, novita, hyperbolic, upstage, moonshot, minimax, siliconflow, replicate, anyscale, cerebras, ollama, deepinfra, sambanova, nebius, baseten, cloudflare, ai21, nvidia, friendli, venice, ovhcloud, scaleway, stepfun, xiaomi, liquid, arcee, chutes, morph, nextbit, parasail, phala, modelrun

```bash
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_OPENAI_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "openai/gpt-4o",
    "messages": [{"role": "user", "content": "Hello, how are you?"}]
  }'
```

#### 2. Azure OpenAI

```bash
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "azure/gpt-4",
    "messages": [{"role": "user", "content": "Hello, how are you?"}],
    "options": {
      "azureEndpoint": "https://your-resource.openai.azure.com"
    }
  }'
```

#### 3. Anthropic

```bash
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_ANTHROPIC_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "anthropic/claude-3-5-sonnet-20241022",
    "messages": [{"role": "user", "content": "Hello, how are you?"}]
  }'
```

#### 4. Google Gemini

```bash
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization:Bearer YOUR_GOOGLE_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gemini/gemini-2.0-flash",
    "messages": [{"role": "user", "content": "Hello, how are you?"}]
  }'
```

#### 5. AWS Bedrock

```bash
# With direct model invocation
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "bedrock/anthropic.claude-3-sonnet-20240229-v1:0",
    "messages": [{"role": "user", "content": "Hello"}],
    "options": {
      "awsAccessKeyID": "AKIAIOSFODNN7EXAMPLE",
      "awsSecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
      "awsRegion": "us-east-1"
    }
  }'

# With Inference Profile (region optional - extracted from ARN)
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "bedrock/anthropic.claude-3-sonnet-20240229-v1:0",
    "messages": [{"role": "user", "content": "Hello"}],
    "options": {
      "awsAccessKeyID": "AKIAIOSFODNN7EXAMPLE",
      "awsSecretAccessKey": "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
      "awsInferenceProfileARN": "arn:aws:bedrock:us-east-1:123456789012:inference-profile/my-profile"
    }
  }'
```

#### 6. Google Vertex AI

```bash
curl -X POST http://localhost:8082/api/v1/chat/completions \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "model": "vertex/gemini-2.0-flash",
    "messages": [{"role": "user", "content": "Hello"}],
    "options": {
      "gcpProjectID": "my-gcp-project-123",
      "gcpLocation": "us-central1"
    }
  }'
```

---

## Provider Testing Status

| **Provider** | **Non-Streaming** | **Streaming** | **Structured Output** | **Reasoning** |
|-------------|:-----------------:|:-------------:|:---------------------:|:-------------:|
| **OpenAI** | ✅ | ✅ | ✅ | ✅ |
| **Azure OpenAI** | ✅ | ✅ | ✅ | ✅ |
| **Anthropic** | ○ | ○ | ○ | ○ |
| **Google Gemini** | ✅ | ✅ | ✅ | ✅ |
| **AWS Bedrock** | ○ | ○ | ○ | ○ |
| **Vertex AI** | ○ | ○ | ○ | ○ |
| **DeepSeek** | ○ | ○ | ○ | ○ |
| **Mistral** | ○ | ○ | ○ | ○ |
| **Groq** | ○ | ○ | ○ | ○ |
| **Together AI** | ○ | ○ | ○ | ○ |
| **Fireworks** | ○ | ○ | ○ | ○ |
| **Perplexity** | ○ | ○ | ○ | ○ |
| **Cohere** | ○ | ○ | ○ | ○ |
| **Novita** | ○ | ○ | ○ | ○ |
| **Hyperbolic** | ○ | ○ | ○ | ○ |
| **Upstage** | ○ | ○ | ○ | ○ |
| **Moonshot** | ○ | ○ | ○ | ○ |
| **MiniMax** | ○ | ○ | ○ | ○ |
| **SiliconFlow** | ○ | ○ | ○ | ○ |
| **Replicate** | ○ | ○ | ○ | ○ |
| **Anyscale** | ○ | ○ | ○ | ○ |
| **Cerebras** | ○ | ○ | ○ | ○ |
| **Ollama** | ○ | ○ | ○ | ○ |
| **DeepInfra** | ○ | ○ | ○ | ○ |
| **SambaNova** | ○ | ○ | ○ | ○ |
| **Nebius** | ○ | ○ | ○ | ○ |
| **Baseten** | ○ | ○ | ○ | ○ |
| **Cloudflare** | ○ | ○ | ○ | ○ |
| **AI21 Labs** | ○ | ○ | ○ | ○ |
| **NVIDIA** | ○ | ○ | ○ | ○ |
| **FriendliAI** | ○ | ○ | ○ | ○ |
| **Venice** | ○ | ○ | ○ | ○ |
| **OVHCloud** | ○ | ○ | ○ | ○ |
| **Scaleway** | ○ | ○ | ○ | ○ |
| **StepFun** | ○ | ○ | ○ | ○ |
| **Xiaomi MiMo** | ○ | ○ | ○ | ○ |
| **Liquid AI** | ○ | ○ | ○ | ○ |
| **Arcee** | ○ | ○ | ○ | ○ |
| **Chutes** | ○ | ○ | ○ | ○ |
| **Morph** | ○ | ○ | ○ | ○ |
| **NextBit** | ○ | ○ | ○ | ○ |
| **Parasail** | ○ | ○ | ○ | ○ |
| **Phala** | ○ | ○ | ○ | ○ |
| **ModelRun** | ○ | ○ | ○ | ○ |
| **xAI** | ○ | ○ | ○ | ○ |
| **Zhipu AI (Zai)** | ✅ | ○ | ○ | ○ |

### Legend

- ✅ **Tested** - Feature verified working
- ○ **Not Tested** - Implementation complete, testing pending
- ❌ **Not Supported** - Feature not available for this provider

---

## License

MIT License
