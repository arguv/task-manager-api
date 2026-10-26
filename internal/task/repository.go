package task

import "sync"

// Repository defines the storage operations required by the service.
type Repository interface {
	List() []Task
	Get(id int) (Task, error)
	Create(task Task) (Task, error)
	Update(task Task) (Task, error)
	Delete(id int) error
}

// MemoryRepository stores tasks in memory.
type MemoryRepository struct {
	mu     sync.Mutex
	tasks  []Task
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks: []Task{
			{ID: 1, Title: "Learn Go", Status: "in_progress", Priority: "high"},
			{ID: 2, Title: "Build REST API", Status: "todo", Priority: "medium"},
		},
		nextID: 3,
	}
}

func (r *MemoryRepository) List() []Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	return append([]Task{}, r.tasks...)
}

func (r *MemoryRepository) Get(id int) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, currentTask := range r.tasks {
		if currentTask.ID == id {
			return currentTask, nil
		}
	}

	return Task{}, ErrNotFound
}

func (r *MemoryRepository) Create(task Task) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	task.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, task)

	return task, nil
}

func (r *MemoryRepository) Update(task Task) (Task, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i := range r.tasks {
		if r.tasks[i].ID == task.ID {
			r.tasks[i] = task
			return task, nil
		}
	}

	return Task{}, ErrNotFound
}

func (r *MemoryRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, currentTask := range r.tasks {
		if currentTask.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return nil
		}
	}

	return ErrNotFound
}

// Compile-time check that MemoryRepository implements Repository.
var _ Repository = (*MemoryRepository)(nil)
