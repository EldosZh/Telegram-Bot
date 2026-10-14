package todo

import (
	"time"
)

// Task represents a single user task in the task manager.
type Task struct {
	ID        int       `json:"id"`
	UserID    int64     `json:"user_id"`
	Title     string    `json:"title"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
}
