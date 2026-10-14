package storage

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"telegram-task-bot/internal/todo"
)

var ErrCorruptedJSON = errors.New("task data file contains invalid JSON")

type JSONStore struct {
	mu   sync.RWMutex
	path string
}

func NewJSONStore(path string) (*JSONStore, error) {
	if path == "" {
		return nil, errors.New("storage path is required")
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("create storage dir: %w", err)
	}

	return &JSONStore{path: path}, nil
}

func (s *JSONStore) Load() ([]todo.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, err := os.Stat(s.path); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, nil
		}
		return nil, fmt.Errorf("read task file: %w", err)
	}

	data, err := os.ReadFile(s.path)
	if err != nil {
		return nil, fmt.Errorf("read task file: %w", err)
	}

	if len(data) == 0 || len(bytes.TrimSpace(data)) == 0 {
		return nil, nil
	}

	var tasks []todo.Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCorruptedJSON, err)
	}

	return tasks, nil
}

func (s *JSONStore) Save(tasks []todo.Task) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := os.MkdirAll(filepath.Dir(s.path), 0o755); err != nil {
		return fmt.Errorf("create storage dir: %w", err)
	}

	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	tmpPath := s.path + ".tmp"
	if err := os.WriteFile(tmpPath, data, 0o644); err != nil {
		return fmt.Errorf("write temp file: %w", err)
	}

	if err := os.Rename(tmpPath, s.path); err != nil {
		_ = os.Remove(tmpPath)
		return fmt.Errorf("move task file into place: %w", err)
	}

	return nil
}
