package bot

import (
	"fmt"
	"strconv"
	"strings"

	"telegram-task-bot/internal/todo"
)

func (b *Bot) addTask(chatID, userID int64, args string) {
	title := strings.TrimSpace(args)
	if title == "" {
		b.sendText(chatID, "❌ Please provide a task.")
		return
	}

	task, err := b.service.Add(userID, title)
	if err != nil {
		if err == todo.ErrEmptyTitle {
			b.sendText(chatID, "❌ Please provide a task.")
			return
		}
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to add task: %v", err)
		}
		return
	}

	if err := b.persist(); err != nil {
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to save tasks: %v", err)
		}
		return
	}

	b.sendText(chatID, fmt.Sprintf("✅ Task added!\n\n%d. %s", task.ID, task.Title))
}

func (b *Bot) listTasks(chatID, userID int64) {
	tasks := b.service.List(userID)
	if len(tasks) == 0 {
		b.sendText(chatID, "📭 You don't have any tasks yet.\n\nAdd one with:\n/add Learn Go")
		return
	}

	b.sendText(chatID, renderTaskList(tasks))
}

func (b *Bot) completeTask(chatID, userID int64, args string) {
	idText := strings.TrimSpace(args)
	if idText == "" {
		b.sendText(chatID, "❌ Please provide a task ID.")
		return
	}

	id, err := strconv.Atoi(idText)
	if err != nil {
		b.sendText(chatID, "❌ Please provide a valid task ID.")
		return
	}

	task, err := b.service.Complete(userID, id)
	if err != nil {
		if err == todo.ErrTaskNotFound {
			b.sendText(chatID, fmt.Sprintf("❌ Task with ID %d not found.", id))
			return
		}
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to complete task: %v", err)
		}
		return
	}

	if err := b.persist(); err != nil {
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to save tasks: %v", err)
		}
		return
	}

	b.sendText(chatID, fmt.Sprintf("✅ Task completed!\n\n%d. %s", task.ID, task.Title))
}

func (b *Bot) deleteTask(chatID, userID int64, args string) {
	idText := strings.TrimSpace(args)
	if idText == "" {
		b.sendText(chatID, "❌ Please provide a task ID.")
		return
	}

	id, err := strconv.Atoi(idText)
	if err != nil {
		b.sendText(chatID, "❌ Please provide a valid task ID.")
		return
	}

	if err := b.service.Delete(userID, id); err != nil {
		if err == todo.ErrTaskNotFound {
			b.sendText(chatID, fmt.Sprintf("❌ Task with ID %d not found.", id))
			return
		}
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to delete task: %v", err)
		}
		return
	}

	if err := b.persist(); err != nil {
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to save tasks: %v", err)
		}
		return
	}

	b.sendText(chatID, "🗑 Task deleted!")
}

func (b *Bot) clearCompleted(chatID, userID int64) {
	removed := b.service.ClearCompleted(userID)
	if err := b.persist(); err != nil {
		b.sendText(chatID, "❌ Something went wrong. Please try again.")
		if b.logger != nil {
			b.logger.Printf("failed to save tasks: %v", err)
		}
		return
	}

	if removed == 0 {
		b.sendText(chatID, "Nothing to clear.")
		return
	}

	b.sendText(chatID, fmt.Sprintf("🧹 Removed %d completed tasks.", removed))
}

func (b *Bot) showStats(chatID, userID int64) {
	stats := b.service.Stats(userID)
	b.sendText(chatID, renderStats(stats))
}

func renderTaskList(tasks []todo.Task) string {
	lines := make([]string, 0, len(tasks)+2)
	lines = append(lines, "📋 Your tasks")
	for _, task := range tasks {
		icon := "⬜"
		if task.Completed {
			icon = "✅"
		}
		lines = append(lines, fmt.Sprintf("%d. %s %s", task.ID, icon, task.Title))
	}
	return strings.Join(lines, "\n")
}

func renderStats(stats todo.Stats) string {
	return fmt.Sprintf("📊 Task Statistics\n\nTotal: %d\nCompleted: %d\nPending: %d\nProgress: %d%%", stats.Total, stats.Completed, stats.Pending, stats.Progress)
}

func helpMessage() string {
	return "📋 Available commands\n\n/add <task> - Add a task\n/list - Show all tasks\n/done <id> - Mark a task as completed\n/delete <id> - Delete a task\n/clear - Delete completed tasks\n/stats - Show statistics\n/help - Show this message"
}

func startMessage() string {
	return "👋 Welcome to Task Manager Bot!\n\nI can help you manage your daily tasks.\n\nCommands:\n\n/add <task> - Add a new task\n/list - Show your tasks\n/done <id> - Complete a task\n/delete <id> - Delete a task\n/clear - Remove completed tasks\n/stats - Show task statistics\n/help - Show available commands\n\nExample:\n\n/add Learn Go"
}
