package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"telegram-task-bot/internal/todo"
)

func TestJSONStoreSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "tasks.json")

	store, err := NewJSONStore(filePath)
	if err != nil {
		t.Fatalf("NewJSONStore returned error: %v", err)
	}

	tasks := []todo.Task{
		{ID: 1, UserID: 42, Title: "Learn Go", Completed: false, CreatedAt: time.Now().UTC()},
		{ID: 2, UserID: 42, Title: "Write tests", Completed: true, CreatedAt: time.Now().UTC()},
	}

	if err := store.Save(tasks); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	loaded, err := store.Load()
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if len(loaded) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(loaded))
	}
	if loaded[0].Title != "Learn Go" {
		t.Fatalf("expected first task title to be Learn Go, got %q", loaded[0].Title)
	}
}

func TestJSONStoreHandlesEmptyAndCorruptedFiles(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "tasks.json")

	store, err := NewJSONStore(filePath)
	if err != nil {
		t.Fatalf("NewJSONStore returned error: %v", err)
	}

	if _, err := store.Load(); err != nil {
		t.Fatalf("Load on missing file should be nil: %v", err)
	}

	if err := os.WriteFile(filePath, []byte(""), 0o644); err != nil {
		t.Fatalf("failed to create empty file: %v", err)
	}
	if tasks, err := store.Load(); err != nil || len(tasks) != 0 {
		t.Fatalf("expected empty file to load as empty tasks, got tasks=%v err=%v", tasks, err)
	}

	if err := os.WriteFile(filePath, []byte("{broken json"), 0o644); err != nil {
		t.Fatalf("failed to write malformed JSON: %v", err)
	}
	if _, err := store.Load(); err == nil {
		t.Fatal("expected corrupted JSON to return an error")
	}
}
