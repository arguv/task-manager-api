package task

import (
	"errors"
	"testing"
)

func TestServiceCreateTask(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	input := CreateTaskRequest{
		Title:    "Learn Go testing",
		Status:   "todo",
		Priority: "high",
	}

	result, err := service.Create(input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != 3 {
		t.Errorf("Expected task ID 3, got %d", result.ID)
	}

	if result.Title != "Learn Go testing" {
		t.Errorf("Unexpected task title: %s", result.Title)
	}

	if result.Status != "todo" {
		t.Errorf("Unexpected task status: %s", result.Status)
	}
}

func TestServiceCreateInvalidTask(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	input := CreateTaskRequest{
		Title:    "",
		Status:   "todo",
		Priority: "high",
	}

	_, err := service.Create(input)

	if !errors.Is(err, ErrInvalid) {
		t.Errorf("Expected ErrInvalid, got: %v", err)
	}
}

func TestServiceGetTask(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	result, err := service.Get(1)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.Title != "Learn Go" {
		t.Errorf("Expected Learn Go, got %s", result.Title)
	}
}

func TestServiceGetTaskNotFound(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	_, err := service.Get(999)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got: %v", err)
	}
}

func TestServiceUpdateTask(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	input := CreateTaskRequest{
		Title:    "Learn Advanced Go",
		Status:   "done",
		Priority: "high",
	}

	result, err := service.Update(1, input)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}

	if result.Title != "Learn Advanced Go" {
		t.Errorf("Unexpected title: %s", result.Title)
	}

	if result.Status != "done" {
		t.Errorf("Expected done, got %s", result.Status)
	}

	stored, err := service.Get(1)
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}

	if stored.Title != input.Title {
		t.Errorf("Stored title was not updated: %s", stored.Title)
	}
}

func TestServiceUpdateTaskNotFound(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	input := CreateTaskRequest{
		Title:    "Missing Task",
		Status:   "todo",
		Priority: "medium",
	}

	_, err := service.Update(999, input)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got: %v", err)
	}
}

func TestServiceDeleteTask(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	err := service.Delete(1)

	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	_, err = service.Get(1)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected deleted task to be missing, got: %v", err)
	}

	tasks := service.List()

	if len(tasks) != 1 {
		t.Errorf("Expected 1 remaining task, got %d", len(tasks))
	}
}

func TestServiceDeleteTaskNotFound(t *testing.T) {
	repo := NewMemoryRepository()
	service := NewService(repo)

	err := service.Delete(999)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got: %v", err)
	}
}
