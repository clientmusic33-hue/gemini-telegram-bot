package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/bot"
	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/config"
	"github.com/clientmusic33-hue/gemini-telegram-bot/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger
	logger := utils.NewLogger(cfg.LogLevel)
	logger.Info("Starting Gemini Telegram Bot",
		"version", "1.0.0",
		"model", cfg.GeminiModel,
		"port", cfg.Port,
	)

	// Create bot
	b, err := bot.New(cfg, logger)
	if err != nil {
		logger.Fatal("Failed to create bot", "error", err)
	}

	// Start bot in goroutine
	go func() {
		if err := b.Start(context.Background()); err != nil {
			logger.Error("Bot error", "error", err)
		}
	}()

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	logger.Info("Received signal", "signal", sig)

	b.Stop()
	logger.Info("Bot stopped gracefully")
}
