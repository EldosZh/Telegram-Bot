package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"telegram-task-bot/internal/bot"
	"telegram-task-bot/internal/storage"
	"telegram-task-bot/internal/todo"
)

func main() {
	logger := log.New(os.Stdout, "telegram-task-bot: ", log.LstdFlags|log.LUTC)

	token, err := loadToken()
	if err != nil {
		logger.Fatal(err)
	}

	dataDir := filepath.Join(".", "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		logger.Fatalf("failed to create data directory: %v", err)
	}

	store, err := storage.NewJSONStore(filepath.Join(dataDir, "tasks.json"))
	if err != nil {
		logger.Fatalf("failed to initialize storage: %v", err)
	}

	tasks, err := store.Load()
	if err != nil {
		logger.Fatalf("failed to load tasks: %v", err)
	}

	service := todo.NewService(tasks)
	telegramBot, err := bot.NewBot(token, service, store, logger)
	if err != nil {
		logger.Fatalf("failed to initialize Telegram bot: %v", err)
	}

	logger.Println("bot started")
	if err := telegramBot.Start(context.Background()); err != nil {
		logger.Fatalf("bot stopped: %v", err)
	}
}

func loadToken() (string, error) {
	if token := strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")); token != "" {
		return token, nil
	}

	filePath := filepath.Join(".", ".env")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", errors.New("TELEGRAM_BOT_TOKEN environment variable is required")
		}
		return "", fmt.Errorf("read .env file: %w", err)
	}

	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		key, value, ok := strings.Cut(trimmed, "=")
		if !ok {
			continue
		}
		if strings.TrimSpace(key) == "TELEGRAM_BOT_TOKEN" {
			return strings.TrimSpace(value), nil
		}
	}

	return "", errors.New("TELEGRAM_BOT_TOKEN environment variable is required")
}
