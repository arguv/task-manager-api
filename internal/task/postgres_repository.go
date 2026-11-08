package task

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// PostgresRepository stores tasks in PostgreSQL.
type PostgresRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *PostgresRepository {
	return &PostgresRepository{
		pool: pool,
	}
}

// List retrieves all tasks.
func (r *PostgresRepository) List() ([]Task, error) {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	rows, err := r.pool.Query(
		ctx,
		`SELECT id, title, status, priority
         FROM tasks
         ORDER BY id`,
	)
	if err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	defer rows.Close()

	tasks := make([]Task, 0)

	for rows.Next() {
		var task Task

		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Status,
			&task.Priority,
		); err != nil {
			return nil, fmt.Errorf("scan task: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate tasks: %w", err)
	}

	return tasks, nil
}

// Get retrieves a task by ID.
func (r *PostgresRepository) Get(id int) (Task, error) {
	var task Task

	err := r.pool.QueryRow(
		context.Background(),
		`SELECT id, title, status, priority
		 FROM tasks
		 WHERE id = $1`,
		id,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Status,
		&task.Priority,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}

	if err != nil {
		return Task{}, fmt.Errorf("get task: %w", err)
	}

	return task, nil
}

// Create inserts a new task.
func (r *PostgresRepository) Create(task Task) (Task, error) {
	err := r.pool.QueryRow(
		context.Background(),
		`INSERT INTO tasks (title, status, priority)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		task.Title,
		task.Status,
		task.Priority,
	).Scan(&task.ID)

	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

// Update modifies an existing task.
func (r *PostgresRepository) Update(task Task) (Task, error) {
	err := r.pool.QueryRow(
		context.Background(),
		`UPDATE tasks
		 SET title = $1,
		     status = $2,
		     priority = $3,
		     updated_at = NOW()
		 WHERE id = $4
		 RETURNING id, title, status, priority`,
		task.Title,
		task.Status,
		task.Priority,
		task.ID,
	).Scan(
		&task.ID,
		&task.Title,
		&task.Status,
		&task.Priority,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return Task{}, ErrNotFound
	}

	if err != nil {
		return Task{}, fmt.Errorf("update task: %w", err)
	}

	return task, nil
}

// Delete removes a task by ID.
func (r *PostgresRepository) Delete(id int) error {
	result, err := r.pool.Exec(
		context.Background(),
		`DELETE FROM tasks WHERE id = $1`,
		id,
	)

	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

// Compile-time check that PostgresRepository implements Repository.
var _ Repository = (*PostgresRepository)(nil)
