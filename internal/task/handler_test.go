package task

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestHandler() *Handler {
	repo := NewMemoryRepository()
	service := NewService(repo)
	return NewHandler(service)
}

func TestHandlerListTasks(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks",
		nil,
	)
	rec := httptest.NewRecorder()

	handler.ListTasks(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", rec.Code)
	}

	var tasks []Task

	if err := json.Unmarshal(rec.Body.Bytes(), &tasks); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestHandlerCreateTask(t *testing.T) {
	handler := newTestHandler()

	body := `{
        "title": "Learn HTTP testing",
        "status": "todo",
        "priority": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected 201, got %d: %s",
			rec.Code, rec.Body.String())
	}

	var created Task

	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if created.ID != 3 {
		t.Errorf("Expected ID 3, got %d", created.ID)
	}

	if created.Title != "Learn HTTP testing" {
		t.Errorf("Unexpected title: %s", created.Title)
	}
}

func TestHandlerCreateInvalidTask(t *testing.T) {
	handler := newTestHandler()

	body := `{
        "title": "",
        "status": "todo",
        "priority": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestHandlerGetTask(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/1",
		nil,
	)

	// Direct handler calls do not automatically populate route parameters.
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.GetTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s",
			rec.Code, rec.Body.String())
	}

	var result Task

	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if result.ID != 1 {
		t.Errorf("Expected ID 1, got %d", result.ID)
	}
}

func TestHandlerUpdateTask(t *testing.T) {
	handler := newTestHandler()

	body := `{
        "title": "Learn Advanced Go",
        "status": "done",
        "priority": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/1",
		strings.NewReader(body),
	)
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.UpdateTask(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d: %s",
			rec.Code, rec.Body.String())
	}

	var updated Task

	if err := json.Unmarshal(rec.Body.Bytes(), &updated); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if updated.ID != 1 {
		t.Errorf("Expected ID 1, got %d", updated.ID)
	}

	if updated.Title != "Learn Advanced Go" {
		t.Errorf("Unexpected title: %s", updated.Title)
	}

	if updated.Status != "done" {
		t.Errorf("Expected done, got %s", updated.Status)
	}
}

func TestHandlerUpdateTaskNotFound(t *testing.T) {
	handler := newTestHandler()

	body := `{
        "title": "Missing Task",
        "status": "todo",
        "priority": "medium"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/999",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.UpdateTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rec.Code)
	}
}

func TestHandlerUpdateInvalidTask(t *testing.T) {
	handler := newTestHandler()

	body := `{
        "title": "",
        "status": "todo",
        "priority": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/1",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.UpdateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestHandlerDeleteTask(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/1",
		nil,
	)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("Expected 204, got %d", rec.Code)
	}

	if rec.Body.Len() != 0 {
		t.Errorf("Expected empty response body, got %s",
			rec.Body.String())
	}

	_, err := handler.service.Get(1)

	if !errors.Is(err, ErrNotFound) {
		t.Errorf("Expected task to be deleted, got: %v", err)
	}
}

func TestHandlerDeleteTaskNotFound(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/999",
		nil,
	)
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rec.Code)
	}
}

func TestHandlerGetInvalidID(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/abc",
		nil,
	)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	handler.GetTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestHandlerGetTaskNotFound(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/999",
		nil,
	)
	req.SetPathValue("id", "999")

	rec := httptest.NewRecorder()

	handler.GetTask(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected 404, got %d", rec.Code)
	}
}

func TestHandlerCreateMalformedJSON(t *testing.T) {
	handler := newTestHandler()

	body := `{"title": "Learn Go",`

	req := httptest.NewRequest(
		http.MethodPost,
		"/api/v1/tasks",
		strings.NewReader(body),
	)

	rec := httptest.NewRecorder()

	handler.CreateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestHandlerUpdateInvalidID(t *testing.T) {
	handler := newTestHandler()

	body := `{
        "title": "Updated Task",
        "status": "done",
        "priority": "high"
    }`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/abc",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	handler.UpdateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestHandlerUpdateMalformedJSON(t *testing.T) {
	handler := newTestHandler()

	body := `{"title": "Updated Task",`

	req := httptest.NewRequest(
		http.MethodPut,
		"/api/v1/tasks/1",
		strings.NewReader(body),
	)
	req.SetPathValue("id", "1")

	rec := httptest.NewRecorder()

	handler.UpdateTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}

func TestHandlerDeleteInvalidID(t *testing.T) {
	handler := newTestHandler()

	req := httptest.NewRequest(
		http.MethodDelete,
		"/api/v1/tasks/abc",
		nil,
	)
	req.SetPathValue("id", "abc")

	rec := httptest.NewRecorder()

	handler.DeleteTask(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", rec.Code)
	}
}
