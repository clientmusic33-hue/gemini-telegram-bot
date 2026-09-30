package bot

import (
	"context"
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/config"
	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/gemini"
	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/storage"
	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/utils"
)

type Bot struct {
	api *tgbotapi.BotAPI
	cfg *config.Config
	geminClient *gemini.Client
	storage storage.Storage
	logger *utils.Logger
	done chan struct{}
}

func New(cfg *config.Config, logger *utils.Logger) (*Bot, error) {
	// Create Telegram API client
	api, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram client: %w", err)
	}

	api.Debug = cfg.LogLevel == "debug"

	// Create Gemini client
	geminClient, err := gemini.New(
		cfg.GeminiAPIKey,
		cfg.GeminiModel,
		cfg.GeminiTimeout,
		cfg.GeminiMaxRetries,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gemini client: %w", err)
	}

	// Create storage
	s := storage.NewMemoryStorage()

	return &Bot{
		api: api,
		cfg: cfg,
		geminClient: geminClient,
		storage: s,
		logger: logger,
		done: make(chan struct{}),
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	b.logger.Info("Bot started", "username", b.api.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = b.cfg.PollingTimeout

	updates := b.api.GetUpdatesChan(u)

	for {
		select {
		case <-b.done:
			b.logger.Info("Bot stopping")
			return nil

		case update := <-updates:
			if update.Message != nil {
				go b.handleMessage(ctx, update.Message)
			}
		}
	}
}

func (b *Bot) Stop() {
	close(b.done)
	b.geminClient.Close()
	b.storage.Close()
}

func (b *Bot) handleMessage(ctx context.Context, msg *tgbotapi.Message) {
	if msg.Text == "" {
		return
	}

	userID := int64(msg.From.ID)
	chatID := msg.Chat.ID

	b.logger.Info("Message received",
		"user_id", userID,
		"username", msg.From.UserName,
		"text", msg.Text,
	)

	// Check if user is blocked
	if b.cfg.IsBlocked(userID) {
		b.sendMessage(chatID, "❌ You are blocked from using this bot.")
		return
	}

	// Check if user is allowed
	if !b.cfg.IsAllowed(userID) {
		b.sendMessage(chatID, "❌ You don't have access to this bot.")
		return
	}

	// Handle commands
	if strings.HasPrefix(msg.Text, "/") {
		b.handleCommand(ctx, chatID, userID, msg.From.UserName, msg.Text)
		return
	}

	// Handle regular messages
	b.handleUserMessage(ctx, chatID, userID, msg.From.UserName, msg.Text)
}

func (b *Bot) handleCommand(ctx context.Context, chatID int64, userID int64, username, text string) {
	command := strings.Fields(text)[0]

	switch command {
	case "/start":
		b.cmdStart(chatID)
	case "/help":
		b.cmdHelp(chatID)
	case "/new":
		b.cmdNewConversation(chatID, userID)
	case "/stats":
		b.cmdStats(chatID, userID)
	case "/model":
		b.cmdModel(chatID)
	case "/reset":
		b.cmdReset(chatID, userID)
	default:
		b.sendMessage(chatID, "Unknown command. Type /help for available commands.")
	}
}

func (b *Bot) handleUserMessage(ctx context.Context, chatID int64, userID int64, username, text string) {
	// Check rate limit
	userData, _ := b.storage.GetUser(userID)
	if userData.MessageCountToday >= int64(b.cfg.DailyMessageLimit) {
		b.sendMessage(chatID, fmt.Sprintf(
			"⚠️ Daily message limit reached (%d messages). Try again tomorrow!",
			b.cfg.DailyMessageLimit,
		))
		return
	}

	// Check input length
	if len(text) > b.cfg.MaxInputCharacters {
		b.sendMessage(chatID, fmt.Sprintf(
			"⚠️ Message too long! Maximum %d characters.",
			b.cfg.MaxInputCharacters,
		))
		return
	}

	// Show typing indicator
	if b.cfg.TypingIndicator {
		b.api.SendChatAction(tgbotapi.NewChatAction(chatID, tgbotapi.ChatTyping))
	}

	// Update message count
	b.storage.UpdateMessageCount(userID)

	// Add user message to history
	b.storage.AddMessage(userID, "user", text)

	// Get conversation history
	userData, _ = b.storage.GetUser(userID)
	var messages []gemini.Message

	for _, msg := range userData.ConversationHistory {
		messages = append(messages, gemini.Message{
			Role: msg.Role,
			Content: msg.Content,
		})
	}

	// Call Gemini API
	resp, err := b.geminClient.SendMessage(ctx, messages)
	if err != nil {
		b.logger.Error("Gemini error", "error", err)

		errMsg := "❌ Error processing your message"
		if strings.Contains(err.Error(), "429") {
			errMsg = "⚠️ Rate limited. Please try again in a moment."
		} else if strings.Contains(err.Error(), "401") {
			errMsg = "⚠️ API authentication failed."
		} else if strings.Contains(err.Error(), "quota") {
			errMsg = "⚠️ API quota exceeded. Try again later."
		}

		b.sendMessage(chatID, errMsg)
		return
	}

	// Add bot response to history
	b.storage.AddMessage(userID, "model", resp)

	// Send response (split if too long)
	b.sendLongMessage(chatID, resp)
}

func (b *Bot) cmdStart(chatID int64) {
	msg := `🤖 **Welcome to Gemini Telegram Bot!**

I'm an AI chatbot powered by Google Gemini API.

✨ **Features:**
• Real-time AI conversations
• Conversation history (last 20 messages)
• Daily message limits
• Admin commands

📝 **Quick Start:**
Just send me any message and I'll respond!

ℹ️ Type /help for all available commands.

🔐 Made with ❤️ for privacy and performance`

	b.sendMessage(chatID, msg)
}

func (b *Bot) cmdHelp(chatID int64) {
	msg := `📚 **Available Commands:**

/start - Welcome message
/help - Show this message
/new - Start a new conversation (clears history)
/stats - Show your usage statistics
/model - Show current AI model
/reset - Reset conversation context

💡 **How to Use:**
1. Send any message to chat with me
2. I'll remember our conversation
3. Use /new to start fresh
4. Use /stats to check your usage

⚠️ **Limits:**
• Daily messages: 50
• Max input: 12,000 characters
• History: 20 messages

🆘 **Need Help?**
Make sure your message:
- Is under 12,000 characters
- Doesn't exceed daily limits
- Uses a clear question format`

	b.sendMessage(chatID, msg)
}

func (b *Bot) cmdNewConversation(chatID int64, userID int64) {
	b.storage.ClearConversation(userID)
	b.sendMessage(chatID, "✅ Conversation cleared! Starting fresh.")
}

func (b *Bot) cmdReset(chatID int64, userID int64) {
	b.storage.ClearConversation(userID)
	b.sendMessage(chatID, "✅ Conversation reset to default state.")
}

func (b *Bot) cmdStats(chatID int64, userID int64) {
	userData, _ := b.storage.GetUser(userID)

	lastMsg := "Never"
	if !userData.LastMessageTime.IsZero() {
		lastMsg = userData.LastMessageTime.Format("2006-01-02 15:04:05")
	}

	msg := fmt.Sprintf(`📊 **Your Statistics:**

Total messages: %d
Messages today: %d / %d
Last message: %s
Conversation messages: %d / %d
Current model: %s
Member since: %s`,
		userData.MessageCount,
		userData.MessageCountToday,
		b.cfg.DailyMessageLimit,
		lastMsg,
		len(userData.ConversationHistory),
		b.cfg.MaxConversationMessages,
		b.cfg.GeminiModel,
		userData.CreatedAt.Format("2006-01-02"),
	)

	b.sendMessage(chatID, msg)
}

func (b *Bot) cmdModel(chatID int64) {
	msg := fmt.Sprintf("🤖 Current model: **%s**", b.cfg.GeminiModel)
	b.sendMessage(chatID, msg)
}

func (b *Bot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdown

	if _, err := b.api.Send(msg); err != nil {
		b.logger.Error("Failed to send message", "error", err)
		// Fallback to plain text
		msg.ParseMode = ""
		b.api.Send(msg)
	}
}

func (b *Bot) sendLongMessage(chatID int64, text string) {
	const maxLen = 4096

	for len(text) > 0 {
		end := maxLen
		if end > len(text) {
			end = len(text)
		}

		// Try to split at newline
		if end < len(text) {
			for i := end; i > 0; i-- {
				if text[i-1] == '\n' {
					end = i
					break
				}
			}
		}

		b.sendMessage(chatID, text[:end])
		text = text[end:]

		if len(text) > 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}
}
