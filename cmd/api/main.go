package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arguv/task-manager-api/internal/task"
)

type HealthResponse struct {
	Status  string `json:"status"`
	Service string `json:"service"`
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

func main() {
	taskRepository := task.NewMemoryRepository()
	taskService := task.NewService(taskRepository)
	taskHandler := task.NewHandler(taskService)

	router := task.NewRouter(taskHandler)
	router.HandleFunc("GET /health", healthCheck)

	log.Println("Server starting on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
