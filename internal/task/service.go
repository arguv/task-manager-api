package task

import (
	"errors"
	"strings"
)

var (
	ErrNotFound = errors.New("task not found")
	ErrInvalid  = errors.New("invalid task")
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func validateTask(input CreateTaskRequest) bool {
	if strings.TrimSpace(input.Title) == "" || len([]rune(strings.TrimSpace(input.Title))) > 255 {
		return false
	}

	switch input.Status {
	case "todo", "in_progress", "done":
	default:
		return false
	}

	switch input.Priority {
	case "low", "medium", "high":
	default:
		return false
	}

	return true
}

func (s *Service) List() ([]Task, error) {
	return s.repo.List()
}

func (s *Service) Get(id int) (Task, error) {
	return s.repo.Get(id)
}

func (s *Service) Create(input CreateTaskRequest) (Task, error) {
	if !validateTask(input) {
		return Task{}, ErrInvalid
	}

	newTask := Task{
		Title:    strings.TrimSpace(input.Title),
		Status:   input.Status,
		Priority: input.Priority,
	}

	return s.repo.Create(newTask)
}

func (s *Service) Update(id int, input CreateTaskRequest) (Task, error) {
	if !validateTask(input) {
		return Task{}, ErrInvalid
	}

	updatedTask := Task{
		ID:       id,
		Title:    strings.TrimSpace(input.Title),
		Status:   input.Status,
		Priority: input.Priority,
	}

	return s.repo.Update(updatedTask)
}

func (s *Service) Delete(id int) error {
	return s.repo.Delete(id)
}
