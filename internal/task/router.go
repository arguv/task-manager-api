package task

import "net/http"

func NewRouter(handler *Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/tasks", handler.ListTasks)
	mux.HandleFunc("GET /api/v1/tasks/{id}", handler.GetTask)
	mux.HandleFunc("POST /api/v1/tasks", handler.CreateTask)
	mux.HandleFunc("PUT /api/v1/tasks/{id}", handler.UpdateTask)
	mux.HandleFunc("DELETE /api/v1/tasks/{id}", handler.DeleteTask)

	return mux
}
