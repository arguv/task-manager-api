package task

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) ListTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.List()
	if err != nil {
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	status, priority := r.URL.Query().Get("status"), r.URL.Query().Get("priority")
	filtered := tasks[:0]
	for _, current := range tasks {
		if (status == "" || current.Status == status) && (priority == "" || current.Priority == priority) {
			filtered = append(filtered, current)
		}
	}
	tasks = filtered
	field := r.URL.Query().Get("sort")
	if field != "" && field != "id" && field != "title" && field != "status" && field != "priority" {
		writeJSONError(w, "Invalid sort field", http.StatusBadRequest)
		return
	}
	if field != "" {
		sort.Slice(tasks, func(i, j int) bool {
			switch field {
			case "id":
				return tasks[i].ID < tasks[j].ID
			case "title":
				return tasks[i].Title < tasks[j].Title
			case "status":
				return tasks[i].Status < tasks[j].Status
			default:
				return tasks[i].Priority < tasks[j].Priority
			}
		})
	}
	page, pageSize, paginationErr := pagination(r)
	if paginationErr != nil {
		writeJSONError(w, paginationErr.Error(), http.StatusBadRequest)
		return
	}
	start := (page - 1) * pageSize
	if start >= len(tasks) {
		tasks = []Task{}
	} else {
		end := start + pageSize
		if end > len(tasks) {
			end = len(tasks)
		}
		tasks = tasks[start:end]
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		// Response may already have been written.
		log.Printf("Failed to encode tasks: %v", err)
	}
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	result, err := h.service.Get(id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSONError(w, "Task not found", http.StatusNotFound)
		} else {
			writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Failed to encode task: %v", err)
	}
}

func writeJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

func pagination(r *http.Request) (int, int, error) {
	page, size := 1, 50
	var err error
	if value := r.URL.Query().Get("page"); value != "" {
		page, err = strconv.Atoi(value)
	}
	if err != nil || page < 1 {
		return 0, 0, errors.New("page must be a positive integer")
	}
	if value := r.URL.Query().Get("page_size"); value != "" {
		size, err = strconv.Atoi(value)
	}
	if err != nil || size < 1 || size > 100 {
		return 0, 0, errors.New("page_size must be between 1 and 100")
	}
	return page, size, nil
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	var input CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Create(input)
	if err != nil {
		if errors.Is(err, ErrInvalid) {
			writeJSONError(w, "Invalid task data", http.StatusBadRequest)
			return
		}
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Failed to encode created task: %v", err)
	}
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	var input CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeJSONError(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	result, err := h.service.Update(id, input)
	if err != nil {
		if errors.Is(err, ErrInvalid) {
			writeJSONError(w, "Invalid task data", http.StatusBadRequest)
			return
		}
		if errors.Is(err, ErrNotFound) {
			writeJSONError(w, "Task not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(result); err != nil {
		log.Printf("Failed to encode updated task: %v", err)
	}
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		writeJSONError(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(id); err != nil {
		if errors.Is(err, ErrNotFound) {
			writeJSONError(w, "Task not found", http.StatusNotFound)
			return
		}
		writeJSONError(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
