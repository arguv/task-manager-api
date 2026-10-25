package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
}

type Task struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

type CreateTaskRequest struct {
	Title    string `json:"title"`
	Status   string `json:"status"`
	Priority string `json:"priority"`
}

var tasks = []Task{
	{ID: 1, Title: "Learn Go", Status: "in_progress", Priority: "high"},
	{ID: 2, Title: "Build REST API", Status: "todo", Priority: "medium"},
}

func validateCreateTask(input CreateTaskRequest) string {
	if strings.TrimSpace(input.Title) == "" {
		return "Title is required"
	}

	switch input.Status {
	case "todo", "in_progress", "done":
	default:
		return "Invalid status"
	}

	switch input.Priority {
	case "low", "medium", "high":
	default:
		return "Invalid priority"
	}

	return ""
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:  "ok",
		Service: "task-manager-api",
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode health response: %v", err)
	}
}

func listTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("Failed to encode tasks: %v", err)
	}
}

func getTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for _, task := range tasks {
		if task.ID == id {
			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(task); err != nil {
				log.Printf("Failed to encode task: %v", err)
			}
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var input CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if message := validateCreateTask(input); message != "" {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	task := Task{
		ID:       len(tasks) + 1,
		Title:    strings.TrimSpace(input.Title),
		Status:   input.Status,
		Priority: input.Priority,
	}

	tasks = append(tasks, task)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("Failed to encode created task: %v", err)
	}
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var input CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if message := validateCreateTask(input); message != "" {
		http.Error(w, message, http.StatusBadRequest)
		return
	}

	for i := range tasks {
		if tasks[i].ID == id {
			tasks[i].Title = strings.TrimSpace(input.Title)
			tasks[i].Status = input.Status
			tasks[i].Priority = input.Priority

			w.Header().Set("Content-Type", "application/json")

			if err := json.NewEncoder(w).Encode(tasks[i]); err != nil {
				log.Printf("Failed to encode updated task: %v", err)
			}
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	for i, task := range tasks {
		if task.ID == id {
			tasks = append(tasks[:i], tasks[i+1:]...)

			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Task not found", http.StatusNotFound)
}

func main() {
	http.HandleFunc("GET /health", healthCheck)
	http.HandleFunc("GET /api/v1/tasks", listTasks)
	http.HandleFunc("GET /api/v1/tasks/{id}", getTask)
	http.HandleFunc("POST /api/v1/tasks", createTask)
	http.HandleFunc("PUT /api/v1/tasks/{id}", updateTask)
	http.HandleFunc("DELETE /api/v1/tasks/{id}", deleteTask)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
