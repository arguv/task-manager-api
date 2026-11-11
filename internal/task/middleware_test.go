package task

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAPIKeyMiddleware(t *testing.T) {
	h := APIKeyMiddleware("secret")(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusNoContent) }))
	for _, test := range []struct {
		name, key string
		status    int
	}{{"missing", "", http.StatusUnauthorized}, {"wrong", "bad", http.StatusUnauthorized}, {"valid", "secret", http.StatusNoContent}} {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", nil)
			req.Header.Set("X-API-Key", test.key)
			rec := httptest.NewRecorder()
			h.ServeHTTP(rec, req)
			if rec.Code != test.status {
				t.Fatalf("expected %d, got %d", test.status, rec.Code)
			}
		})
	}
}

func TestRecoverMiddleware(t *testing.T) {
	h := RecoverMiddleware(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { panic("boom") }))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/", nil))
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", rec.Code)
	}
}

func TestListPaginationAndFilters(t *testing.T) {
	h := NewHandler(NewService(NewMemoryRepository()))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks?status=todo&page=1&page_size=1&sort=title", nil)
	rec := httptest.NewRecorder()
	h.ListTasks(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}
	if rec.Body.String() == "[]\n" {
		t.Fatal("expected filtered task")
	}
}
