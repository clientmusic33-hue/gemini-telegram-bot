# Gemini Telegram Bot

🤖 **Production-Ready Telegram AI Bot** powered by Google Gemini API

[![Go Version](https://img.shields.io/badge/Go-1.22-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green.svg)](LICENSE)
[![Build Status](https://github.com/clientmusic33-hue/gemini-telegram-bot/actions/workflows/ci.yml/badge.svg)](https://github.com/clientmusic33-hue/gemini-telegram-bot/actions)

## 🌟 Features

- ✅ **Direct Gemini API Integration** - No OpenRouter, pure Google Gemini API
- ✅ **Long Polling** - Reliable, no public webhook URL required
- ✅ **Conversation Memory** - Per-user conversation history with configurable limits
- ✅ **Advanced Commands** - `/start`, `/help`, `/model`, `/models`, `/new`, `/reset`, `/stats`, `/settings`
- ✅ **User Management** - Admin system, user allowlisting, blocking
- ✅ **Rate Limiting** - Daily message limits, max input characters
- ✅ **Error Handling** - Graceful Gemini API error handling with exponential backoff
- ✅ **Message Splitting** - Automatic split for Telegram's 4096 char limit
- ✅ **Typing Indicator** - Shows bot is thinking
- ✅ **Health Endpoint** - Render compatibility with `/health` check
- ✅ **Structured Logging** - Production-grade logging
- ✅ **Graceful Shutdown** - Handles SIGINT/SIGTERM
- ✅ **Docker Ready** - Multi-stage Dockerfile for production
- ✅ **Render Deployment** - One-click deployment with `render.yaml`
- ✅ **Unit Tests** - Comprehensive test coverage
- ✅ **Free to Use** - MIT license, open source

## 🚀 Quick Start

### Prerequisites

- Go 1.22+
- Telegram account
- Google Cloud account with Gemini API access

### 1. Get Gemini API Key

1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click "Get API Key"
3. Create a new API key in your default project
4. Copy the key (keep it secret!)

### 2. Create Telegram Bot

1. Open Telegram and search for `@BotFather`
2. Send `/newbot`
3. Follow prompts to create your bot
4. Copy the bot token (format: `123456789:ABCDefGHIjklmnoPQRstUvWxyz`)

### 3. Clone & Setup Locally

```bash
git clone https://github.com/clientmusic33-hue/gemini-telegram-bot.git
cd gemini-telegram-bot

cp .env.example .env
# Edit .env and add your credentials
# TELEGRAM_BOT_TOKEN=your_token
# GEMINI_API_KEY=your_key

go mod download
go run ./cmd/bot
```

### 4. Test Your Bot

In Telegram:
1. Search for your bot username
2. Click START or send `/start`
3. Send a message like "Hello!" or "Who is Albert Einstein?"
4. Bot should respond with AI-generated answer

## 📋 Environment Variables

### Required

| Variable | Description | Example |
|----------|-------------|----------|
| `TELEGRAM_BOT_TOKEN` | Telegram bot token from BotFather | `123456789:ABCDefGHI...` |
| `GEMINI_API_KEY` | Google Gemini API key | `AIzaSyD...` |

### Optional

| Variable | Description | Default |
|----------|-------------|----------|
| `GEMINI_MODEL` | Gemini model to use | `gemini-2.0-flash` |
| `ADMIN_IDS` | Comma-separated admin Telegram IDs | (empty) |
| `ALLOWED_USER_IDS` | Comma-separated allowed user IDs (empty=all) | (empty) |
| `BLOCKED_USER_IDS` | Comma-separated blocked user IDs | (empty) |
| `DAILY_MESSAGE_LIMIT` | Max messages per user per day | `50` |
| `MAX_INPUT_CHARACTERS` | Max characters per message | `12000` |
| `MAX_CONVERSATION_MESSAGES` | Max messages in conversation history | `20` |
| `PORT` | HTTP server port | `8080` |
| `LOG_LEVEL` | Logging level (debug/info/warn/error) | `info` |
| `REQUEST_TIMEOUT` | Gemini request timeout (seconds) | `30` |
| `POLLING_TIMEOUT` | Telegram polling timeout (seconds) | `30` |
| `STORAGE_TYPE` | Storage type (memory/sqlite) | `memory` |

## 📁 Project Structure

```
gemini-telegram-bot/
├── cmd/
│   └── bot/
│       └── main.go              # Application entry point
├── internal/
│   ├── bot/
│   │   ├── bot.go              # Bot core logic
│   │   ├── handlers.go         # Message & command handlers
│   │   ├── commands.go         # Command implementations
│   │   ├── middleware.go       # Auth & rate limiting
│   │   └── bot_test.go         # Tests
│   ├── gemini/
│   │   ├── client.go           # Gemini API client
│   │   ├── models.go           # Data structures
│   │   └── client_test.go      # Tests
│   ├── config/
│   │   ├── config.go           # Configuration loading
│   │   └── config_test.go      # Tests
│   ├── storage/
│   │   ├── storage.go          # User data storage
│   │   └── storage_test.go     # Tests
│   └── utils/
│       ├── utils.go            # Helper functions
│       ├── logger.go           # Structured logging
│       └── utils_test.go       # Tests
├── configs/
│   └── config.example.yaml     # Example config
├── .github/
│   └── workflows/
│       └── ci.yml              # GitHub Actions CI
├── .env.example                # Environment variables template
├── .gitignore                  # Git ignore rules
├── Dockerfile                  # Production Docker image
├── render.yaml                 # Render.com deployment config
├── go.mod                      # Go module definition
├── go.sum                      # Go module checksums
├── README.md                   # This file
└── LICENSE                     # MIT License
```

## 🎮 Available Commands

| Command | Description |
|---------|-------------|
| `/start` | Welcome message and bot info |
| `/help` | Show all commands and how to use them |
| `/model` | Show currently selected AI model |
| `/models` | List available Gemini models |
| `/new` | Start a new conversation (clear history) |
| `/reset` | Reset conversation to default |
| `/stats` | Show your usage statistics |
| `/settings` | Show your account settings |

## 🚀 Deploy to Render

### Method 1: Using render.yaml (Recommended)

1. Push code to GitHub:
```bash
git add .
git commit -m "Initial commit"
git push origin main
```

2. Go to [Render Dashboard](https://dashboard.render.com)
3. Click "Create +" → "Web Service"
4. Connect your GitHub account
5. Select `gemini-telegram-bot` repository
6. Render will auto-detect `render.yaml`
7. Add environment variables:
   - `TELEGRAM_BOT_TOKEN`
   - `GEMINI_API_KEY`
   - `GEMINI_MODEL` = `gemini-2.0-flash`
8. Click "Create Web Service"
9. Wait for deployment (2-3 minutes)
10. Test bot in Telegram

### Method 2: Manual Configuration

1. Create new Web Service on Render
2. Build command: `go build -o bot ./cmd/bot`
3. Start command: `./bot`
4. Instance type: Standard (minimum)
5. Environment: Set all vars from `.env.example`
6. Deploy

### Verify Deployment

```bash
# Check bot responds
# Send /start in Telegram

# Check health endpoint
curl https://your-service.onrender.com/health
# Should return: {"status":"ok"}
```

## 🔧 Configuration File (config.example.yaml)

```yaml
bot:
  name: "Gemini Telegram Bot"
  description: "AI-powered Telegram bot"
  version: "1.0.0"

gemini:
  timeout_seconds: 30
  max_retries: 3
  retry_backoff_base: 2

conversation:
  enabled: true
  max_messages: 20
  cleanup_interval_minutes: 60

limits:
  enabled: true
  daily_messages: 50
  max_input_characters: 12000
  rate_limit_window_minutes: 1
  messages_per_window: 10

users:
  admin_ids: []
  allowed_user_ids: []
  blocked_user_ids: []
  track_usage: true

logging:
  level: info
  format: json
  file: ""

server:
  host: 0.0.0.0
  port: 8080
  health_path: /health
  readiness_path: /ready
```

## 🛠️ Development

### Build from source

```bash
go build -o bin/bot ./cmd/bot
```

### Run tests

```bash
go test ./...
go test -v ./...
go test -cover ./...
```

### Code quality

```bash
go fmt ./...
go vet ./...
golangci-lint run ./...
```

### Docker build

```bash
docker build -t gemini-telegram-bot:latest .
docker run --env-file .env gemini-telegram-bot:latest
```

## 🐛 Troubleshooting

### Bot not responding in Telegram

**Issue:** Bot doesn't reply to messages

**Solutions:**
1. Check bot token is correct: `TELEGRAM_BOT_TOKEN=123456789:ABCDef...`
2. Check bot is running: `ps aux | grep bot`
3. Verify Render logs: Dashboard → Logs tab
4. Test locally first: `go run ./cmd/bot`
5. Check network: `curl https://api.telegram.org/bot<TOKEN>/getMe`

### Invalid Telegram Token

**Error:** `Unauthorized` or `401`

**Solutions:**
1. Copy token exactly from BotFather (no extra spaces)
2. Verify format: `123456789:ABCDef...`
3. Generate new token from BotFather if lost
4. Check `.env` file has correct token

### Invalid Gemini API Key

**Error:** `401 Unauthorized` or `"Invalid API Key"`

**Solutions:**
1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Generate new API key
3. Verify no extra spaces in key
4. Key should start with `AIzaSy...`
5. Check project has Generative Language API enabled

### Gemini 429 (Rate Limited)

**Error:** `Too Many Requests`

**Solutions:**
1. Wait 60+ seconds before next request
2. Bot has exponential backoff built-in
3. Check daily quota not exceeded
4. Upgrade Gemini plan if needed
5. Reduce `DAILY_MESSAGE_LIMIT` in config

### Gemini Quota Exceeded

**Error:** `"Quota exceeded"`

**Solutions:**
1. Bot sends user-friendly message: "⚠️ API quota unavailable. Try later."
2. Check Gemini API usage: [Google Cloud Console](https://console.cloud.google.com)
3. Quota resets daily at midnight UTC
4. Upgrade to paid plan for higher limits
5. Check if billing is enabled on Google Cloud

### Model Not Found

**Error:** `Model not found: gemini-xyz`

**Solutions:**
1. Use valid model: `gemini-2.0-flash` or `gemini-1.5-pro`
2. Check Gemini API docs for available models
3. Set `GEMINI_MODEL=gemini-2.0-flash` (default)
4. Verify model is available in your region/plan

### Render Deployment Failure

**Error:** Build fails or service doesn't start

**Solutions:**
1. Check build logs in Render dashboard
2. Verify Go version >= 1.22: Update `go.mod`
3. All env vars set: `TELEGRAM_BOT_TOKEN`, `GEMINI_API_KEY`
4. Check `render.yaml` syntax
5. Clear build cache: Delete service and recreate
6. Try manual setup: Create Web Service without `render.yaml`

### Go Build Failure

**Error:** `go: command not found` or build errors

**Solutions:**
1. Install Go 1.22+: [golang.org](https://golang.org/dl)
2. Add Go to PATH: `export PATH=$PATH:/usr/local/go/bin`
3. Run `go mod tidy` to fix dependencies
4. Delete `go.sum` and run `go mod download`
5. Check `go.mod` Go version matches installed version

### Bot Repeatedly Restarting

**Symptom:** Bot starts/stops in loop

**Solutions:**
1. Check logs for panic: `Render dashboard → Logs`
2. Verify env vars not empty
3. Check Telegram token validity
4. Increase `POLLING_TIMEOUT` if network issues
5. Add more memory if OOM: Render → Instance Type
6. Check for race conditions: `go test -race ./...`

### Telegram Markdown Errors

**Issue:** Messages don't format correctly

**Solutions:**
1. Bot auto-falls back to plain text on format error
2. Check special chars: `_*[]()`
3. Escape reserved chars: `_` → `\_`
4. Use code blocks: ` ```code``` `
5. Test markdown: [Telegram Bot API Docs](https://core.telegram.org/bots/api#formatting-options)

## 🔐 Security

- **Never commit** `.env` or `config.yaml` with real secrets
- **Always use** `TELEGRAM_BOT_TOKEN` and `GEMINI_API_KEY` environment variables
- **API keys never logged** or exposed in error messages
- **User data isolated** per Telegram user
- **Admin functions** protected with `ADMIN_IDS` verification
- **Input validation** on all user messages
- **Rate limiting** prevents abuse
- **Graceful errors** without exposing internals

## 📊 Performance

- **Sub-second** response times (typical: 500-2000ms for Gemini)
- **Memory efficient** conversation storage
- **Concurrent requests** supported
- **Automatic retries** with exponential backoff
- **Connection pooling** for HTTP requests
- **Render optimized** - works on free tier

## 📈 Usage Statistics

Use `/stats` command to see:
- Total messages sent
- Messages today
- Last message time
- Current model
- Conversation history size

## 🤝 Contributing

Contributions welcome! Areas for enhancement:
- [ ] SQLite persistent storage
- [ ] Streaming responses
- [ ] File/image handling
- [ ] Custom AI personas
- [ ] Advanced metrics/analytics
- [ ] Redis session storage
- [ ] Webhook support

## 📝 License

MIT License - free to use, modify, and distribute.

## 🆘 Support

**Issues or questions?**
1. Check [Troubleshooting](#-troubleshooting) section
2. Search existing GitHub issues
3. Create new issue with detailed description
4. Include logs (without secrets!)

## 🎯 Roadmap

- [x] Core Telegram bot
- [x] Gemini API integration
- [x] Conversation memory
- [x] User management
- [x] Rate limiting
- [x] Error handling
- [x] Render deployment
- [ ] Web dashboard
- [ ] Advanced analytics
- [ ] Multi-model support
- [ ] Voice message support

## 📚 Resources

- [Telegram Bot API](https://core.telegram.org/bots/api)
- [Google Gemini API](https://ai.google.dev/tutorials/go_quickstart)
- [Go Telegram Bot Library](https://github.com/go-telegram-bot-api/telegram-bot-api)
- [Render Documentation](https://render.com/docs)
- [Go Best Practices](https://golang.org/doc/effective_go)

---

**Made with ❤️ for AI enthusiasts**

Star ⭐ this repo if you find it useful!
