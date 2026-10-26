package task

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestRouter() *http.ServeMux {
	repo := NewMemoryRepository()
	service := NewService(repo)
	handler := NewHandler(service)

	return NewRouter(handler)
}

func TestRouterGetTask(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/1",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("Expected 200, got %d", rec.Code)
	}
}

func TestRouterTaskNotFound(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/tasks/999",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected 404, got %d", rec.Code)
	}
}

func TestRouterInvalidMethod(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(
		http.MethodPatch,
		"/api/v1/tasks/1",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("Expected 405, got %d", rec.Code)
	}
}

func TestRouterUnknownEndpoint(t *testing.T) {
	router := newTestRouter()

	req := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/unknown",
		nil,
	)

	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected 404, got %d", rec.Code)
	}
}
