package todo

import (
	"errors"
	"testing"
	"time"
)

func TestAddTask(t *testing.T) {
	service := NewService(nil)

	task, err := service.Add(101, " Learn Go ")
	if err != nil {
		t.Fatalf("Add returned error: %v", err)
	}
	if task.ID != 1 {
		t.Fatalf("expected task ID 1, got %d", task.ID)
	}
	if task.UserID != 101 {
		t.Fatalf("expected user ID 101, got %d", task.UserID)
	}
	if task.Title != "Learn Go" {
		t.Fatalf("expected trimmed title, got %q", task.Title)
	}
	if task.Completed {
		t.Fatal("new task should not be marked completed")
	}
	if task.CreatedAt.IsZero() {
		t.Fatal("created_at should be set")
	}
}

func TestEmptyTitleValidation(t *testing.T) {
	service := NewService(nil)

	_, err := service.Add(1, "   ")
	if !errors.Is(err, ErrEmptyTitle) {
		t.Fatalf("expected ErrEmptyTitle, got %v", err)
	}
}

func TestTaskIDGenerationPerUser(t *testing.T) {
	service := NewService(nil)

	_, _ = service.Add(1, "Task 1")
	_, _ = service.Add(1, "Task 2")
	_, _ = service.Add(2, "Other user task")

	user1Tasks := service.List(1)
	if len(user1Tasks) != 2 {
		t.Fatalf("expected 2 tasks for user 1, got %d", len(user1Tasks))
	}
	if user1Tasks[1].ID != 2 {
		t.Fatalf("expected second task ID to be 2, got %d", user1Tasks[1].ID)
	}

	user2Tasks := service.List(2)
	if len(user2Tasks) != 1 || user2Tasks[0].ID != 1 {
		t.Fatalf("expected user 2 task ID 1, got %#v", user2Tasks)
	}
}

func TestListOnlyCurrentUserTasks(t *testing.T) {
	service := NewService(nil)
	_, _ = service.Add(10, "A")
	_, _ = service.Add(20, "B")
	_, _ = service.Add(10, "C")

	user10 := service.List(10)
	if len(user10) != 2 {
		t.Fatalf("expected 2 tasks for user 10, got %d", len(user10))
	}
	for _, task := range user10 {
		if task.UserID != 10 {
			t.Fatalf("task belongs to another user: %+v", task)
		}
	}
}

func TestCompleteTask(t *testing.T) {
	service := NewService(nil)
	task, _ := service.Add(11, "Complete me")

	completed, err := service.Complete(11, task.ID)
	if err != nil {
		t.Fatalf("Complete returned error: %v", err)
	}
	if !completed.Completed {
		t.Fatal("task should be marked completed")
	}
}

func TestCompleteNonExistentTask(t *testing.T) {
	service := NewService(nil)
	_, err := service.Complete(11, 99)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestUserCannotCompleteAnotherUsersTask(t *testing.T) {
	service := NewService(nil)
	task, _ := service.Add(5, "Owned by user 5")

	_, err := service.Complete(7, task.ID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected task not found for different user, got %v", err)
	}
}

func TestDeleteTask(t *testing.T) {
	service := NewService(nil)
	task, _ := service.Add(12, "Delete me")

	if err := service.Delete(12, task.ID); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
	if len(service.List(12)) != 0 {
		t.Fatal("task should be removed")
	}
}

func TestUserCannotDeleteAnotherUsersTask(t *testing.T) {
	service := NewService(nil)
	task, _ := service.Add(8, "Owned by user 8")

	err := service.Delete(3, task.ID)
	if !errors.Is(err, ErrTaskNotFound) {
		t.Fatalf("expected ErrTaskNotFound, got %v", err)
	}
}

func TestClearCompletedTasks(t *testing.T) {
	service := NewService(nil)
	_, _ = service.Add(1, "Pending")
	completed, _ := service.Add(1, "Completed")
	_, _ = service.Add(2, "Other user")
	_, _ = service.Complete(1, completed.ID)
	_, _ = service.Add(1, "Another completed")
	_, _ = service.Complete(1, 3)

	removed := service.ClearCompleted(1)
	if removed != 2 {
		t.Fatalf("expected 2 removed tasks, got %d", removed)
	}
	if len(service.List(1)) != 1 {
		t.Fatalf("expected 1 remaining task, got %d", len(service.List(1)))
	}
}

func TestStatsCalculation(t *testing.T) {
	service := NewService(nil)
	_, _ = service.Add(1, "One")
	_, _ = service.Add(1, "Two")
	_, _ = service.Add(1, "Three")
	_, _ = service.Add(1, "Four")
	_, _ = service.Add(2, "Other")
	_, _ = service.Complete(1, 1)
	_, _ = service.Complete(1, 3)

	stats := service.Stats(1)
	if stats.Total != 4 {
		t.Fatalf("expected 4 total tasks, got %d", stats.Total)
	}
	if stats.Completed != 2 {
		t.Fatalf("expected 2 completed tasks, got %d", stats.Completed)
	}
	if stats.Pending != 2 {
		t.Fatalf("expected 2 pending tasks, got %d", stats.Pending)
	}
	if stats.Progress != 50 {
		t.Fatalf("expected 50%% progress, got %d%%", stats.Progress)
	}
}

func TestZeroTaskStats(t *testing.T) {
	service := NewService(nil)
	stats := service.Stats(99)
	if stats.Total != 0 || stats.Completed != 0 || stats.Pending != 0 || stats.Progress != 0 {
		t.Fatalf("unexpected zero-task stats: %+v", stats)
	}
}

func TestServiceCanBeCreatedWithExistingTasks(t *testing.T) {
	initial := []Task{{ID: 1, UserID: 1, Title: "Persisted", Completed: false, CreatedAt: time.Now().UTC()}}
	service := NewService(initial)
	if got := len(service.List(1)); got != 1 {
		t.Fatalf("expected existing task to be loaded, got %d", got)
	}
}
