package task

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerRepositoryFailuresReturn500(t *testing.T) {
	tests := []struct {
		name   string
		call   func(*Handler, http.ResponseWriter, *http.Request)
		method string
		path   string
		body   string
	}{
		{"list", (*Handler).ListTasks, http.MethodGet, "/api/v1/tasks", ""},
		{"get", (*Handler).GetTask, http.MethodGet, "/api/v1/tasks/1", ""},
		{"create", (*Handler).CreateTask, http.MethodPost, "/api/v1/tasks", `{"title":"x","status":"todo","priority":"low"}`},
		{"update", (*Handler).UpdateTask, http.MethodPut, "/api/v1/tasks/1", `{"title":"x","status":"todo","priority":"low"}`},
		{"delete", (*Handler).DeleteTask, http.MethodDelete, "/api/v1/tasks/1", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newFailingHandler()
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			req.SetPathValue("id", "1")
			rec := httptest.NewRecorder()
			tt.call(h, rec, req)
			if rec.Code != http.StatusInternalServerError {
				t.Fatalf("expected 500, got %d: %s", rec.Code, rec.Body.String())
			}
		})
	}
}
