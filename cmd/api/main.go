package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"github.com/arguv/task-manager-api/internal/config"
	"github.com/arguv/task-manager-api/internal/database"
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
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; using existing environment variables")
	}

	appConfig, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pool, err := database.NewPostgresPool(appConfig.DatabaseURL)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer pool.Close()

	log.Println("Connected to PostgreSQL successfully")

	taskRepository := task.NewPostgresRepository(pool)
	taskService := task.NewService(taskRepository)
	taskHandler := task.NewHandler(taskService)

	router := task.NewRouter(taskHandler)
	router.HandleFunc("GET /health", healthCheck)

	protectedRouter := task.APIKeyMiddleware(appConfig.APIKey)(router)
	server := newHTTPServer(appConfig.HTTPAddr, task.RecoverMiddleware(task.LoggingMiddleware(protectedRouter)))

	serverErrors := make(chan error, 1)

	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	select {
	case <-stop:
		log.Println("Shutdown signal received")
	case err := <-serverErrors:
		log.Fatalf("HTTP server failed: %v", err)
	}

	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("HTTP shutdown failed: %v", err)
	}

	log.Println("Application stopped")
}
