package task

import (
	"net/http/httptest"
	"testing"
)

func BenchmarkListTasks(b *testing.B) {
	h := NewHandler(NewService(NewMemoryRepository()))
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("GET", "/api/v1/tasks?page_size=10&sort=title", nil)
		h.ListTasks(httptest.NewRecorder(), req)
	}
}
