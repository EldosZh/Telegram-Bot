package bot

import (
	"context"
	"errors"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-task-bot/internal/storage"
	"telegram-task-bot/internal/todo"
)

// Bot wraps the Telegram API client and business logic.
type Bot struct {
	api     *tgbotapi.BotAPI
	service *todo.Service
	store   *storage.JSONStore
	logger  *log.Logger
}

func NewBot(token string, service *todo.Service, store *storage.JSONStore, logger *log.Logger) (*Bot, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("TELEGRAM_BOT_TOKEN environment variable is required")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		api:     api,
		service: service,
		store:   store,
		logger:  logger,
	}, nil
}

func (b *Bot) Start(ctx context.Context) error {
	updates := b.api.GetUpdatesChan(tgbotapi.UpdateConfig{
		Timeout:        60,
		AllowedUpdates: []string{"message"},
	})

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			b.handleUpdate(update)
		}
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {
	if update.Message == nil {
		return
	}

	chatID := update.Message.Chat.ID
	userID := update.Message.From.ID
	text := strings.TrimSpace(update.Message.Text)
	if text == "" {
		b.sendText(chatID, "Use /help to see available commands.")
		return
	}

	if !strings.HasPrefix(text, "/") {
		b.sendText(chatID, "Use /help to see available commands.")
		return
	}

	command, args, _ := strings.Cut(text, " ")
	command = strings.ToLower(command)
	b.handleCommand(chatID, userID, command, args)
}

func (b *Bot) handleCommand(chatID, userID int64, command, args string) {
	switch command {
	case "/start":
		b.sendText(chatID, startMessage())
	case "/help":
		b.sendText(chatID, helpMessage())
	case "/add":
		b.addTask(chatID, userID, args)
	case "/list":
		b.listTasks(chatID, userID)
	case "/done":
		b.completeTask(chatID, userID, args)
	case "/delete":
		b.deleteTask(chatID, userID, args)
	case "/clear":
		b.clearCompleted(chatID, userID)
	case "/stats":
		b.showStats(chatID, userID)
	default:
		b.sendText(chatID, "❓ Unknown command.\n\n/help - Show available commands")
	}
}

func (b *Bot) sendText(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := b.api.Send(msg); err != nil {
		if b.logger != nil {
			b.logger.Printf("failed to send Telegram message: %v", err)
		}
	}
}

func (b *Bot) persist() error {
	if b.store == nil {
		return nil
	}
	return b.store.Save(b.service.AllTasks())
}
