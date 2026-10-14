package todo

import (
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrEmptyTitle   = errors.New("task title cannot be empty")
	ErrTaskNotFound = errors.New("task not found")
)

// Stats contains counts for a user's task collection.
type Stats struct {
	Total     int
	Completed int
	Pending   int
	Progress  int
}

// Service manages tasks for multiple users and provides concurrency-safe access.
type Service struct {
	mu    sync.RWMutex
	tasks []Task
}

func NewService(initial []Task) *Service {
	clone := make([]Task, len(initial))
	copy(clone, initial)
	return &Service{tasks: clone}
}

// AllTasks returns a copy of all tasks in insertion order.
func (s *Service) AllTasks() []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	out := make([]Task, len(s.tasks))
	copy(out, s.tasks)
	return out
}

// Add creates a new task for the given user.
func (s *Service) Add(userID int64, title string) (Task, error) {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return Task{}, ErrEmptyTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	id := s.nextIDForUser(userID)
	task := Task{
		ID:        id,
		UserID:    userID,
		Title:     trimmed,
		Completed: false,
		CreatedAt: time.Now().UTC(),
	}

	s.tasks = append(s.tasks, task)
	return task, nil
}

// List returns all tasks for a single user in insertion order.
func (s *Service) List(userID int64) []Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		if task.UserID == userID {
			result = append(result, task)
		}
	}
	return result
}

// Complete marks a task as completed when it belongs to the user.
func (s *Service) Complete(userID int64, id int) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.tasks {
		if s.tasks[i].UserID == userID && s.tasks[i].ID == id {
			s.tasks[i].Completed = true
			return s.tasks[i], nil
		}
	}

	return Task{}, ErrTaskNotFound
}

// Delete removes a task belonging to the user.
func (s *Service) Delete(userID int64, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, task := range s.tasks {
		if task.UserID == userID && task.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			return nil
		}
	}

	return ErrTaskNotFound
}

// ClearCompleted removes all completed tasks for the user and returns the count removed.
func (s *Service) ClearCompleted(userID int64) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	filtered := make([]Task, 0, len(s.tasks))
	removed := 0
	for _, task := range s.tasks {
		if task.UserID == userID && task.Completed {
			removed++
			continue
		}
		filtered = append(filtered, task)
	}
	if removed > 0 {
		s.tasks = filtered
	}
	return removed
}

// Stats calculates total, completed, pending and progress percentages for a user.
func (s *Service) Stats(userID int64) Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := Stats{}
	for _, task := range s.tasks {
		if task.UserID != userID {
			continue
		}
		stats.Total++
		if task.Completed {
			stats.Completed++
		} else {
			stats.Pending++
		}
	}

	if stats.Total == 0 {
		stats.Progress = 0
		return stats
	}

	stats.Progress = (stats.Completed * 100) / stats.Total
	return stats
}

func (s *Service) nextIDForUser(userID int64) int {
	maxID := 0
	for _, task := range s.tasks {
		if task.UserID == userID && task.ID > maxID {
			maxID = task.ID
		}
	}
	return maxID + 1
}
