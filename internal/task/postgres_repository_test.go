package task

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func setupPostgresTest(t *testing.T) (*PostgresRepository, func()) {
	t.Helper()

	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		t.Fatalf("parse test database config: %v", err)
	}

	// Safety check: never run these tests against the development database.
	if config.ConnConfig.Database != "task_manager_test" {
		t.Fatal("integration tests require database task_manager_test")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		t.Fatalf("create test pool: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Fatalf("connect to test database: %v", err)
	}

	repository := NewPostgresRepository(pool)

	// Use a unique prefix so we only clean up our own test data.
	prefix := "integration_test_" + strings.ReplaceAll(
		t.Name(), "/", "_",
	) + "_" + time.Now().Format("20060102150405.000000000")

	cleanup := func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cleanupCancel()

		_, err := pool.Exec(
			cleanupCtx,
			"DELETE FROM tasks WHERE title LIKE $1",
			prefix+"%",
		)
		if err != nil {
			t.Errorf("cleanup test tasks: %v", err)
		}

		pool.Close()
	}

	// Each test uses the prefix when creating records.
	t.Setenv("CURRENT_TEST_PREFIX", prefix)

	return repository, cleanup
}

func testTaskTitle(suffix string) string {
	return os.Getenv("CURRENT_TEST_PREFIX") + "_" + suffix
}

func TestPostgresRepositoryCreateAndGet(t *testing.T) {
	repository, cleanup := setupPostgresTest(t)
	defer cleanup()

	created, err := repository.Create(Task{
		Title:    testTaskTitle("create"),
		Status:   "todo",
		Priority: "high",
	})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if created.ID <= 0 {
		t.Fatalf("expected generated ID, got %d", created.ID)
	}

	found, err := repository.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if found.ID != created.ID || found.Title != created.Title {
		t.Fatalf("unexpected task: %+v", found)
	}
}

func TestPostgresRepositoryList(t *testing.T) {
	repository, cleanup := setupPostgresTest(t)
	defer cleanup()

	created, err := repository.Create(Task{
		Title:    testTaskTitle("list"),
		Status:   "todo",
		Priority: "medium",
	})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	tasks, err := repository.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	found := false

	for _, current := range tasks {
		if current.ID == created.ID {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("created task %d not found in List()", created.ID)
	}
}

func TestPostgresRepositoryUpdate(t *testing.T) {
	repository, cleanup := setupPostgresTest(t)
	defer cleanup()

	created, err := repository.Create(Task{
		Title:    testTaskTitle("update"),
		Status:   "todo",
		Priority: "low",
	})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	created.Status = "done"
	created.Priority = "high"

	updated, err := repository.Update(created)
	if err != nil {
		t.Fatalf("Update() failed: %v", err)
	}

	if updated.Status != "done" || updated.Priority != "high" {
		t.Fatalf("unexpected updated task: %+v", updated)
	}

	found, err := repository.Get(created.ID)
	if err != nil {
		t.Fatalf("Get() after Update() failed: %v", err)
	}

	if found.Status != "done" {
		t.Fatalf("expected persisted status done, got %s", found.Status)
	}
}

func TestPostgresRepositoryDelete(t *testing.T) {
	repository, cleanup := setupPostgresTest(t)
	defer cleanup()

	created, err := repository.Create(Task{
		Title:    testTaskTitle("delete"),
		Status:   "todo",
		Priority: "medium",
	})
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if err := repository.Delete(created.ID); err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	_, err = repository.Get(created.ID)

	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestPostgresRepositoryNotFound(t *testing.T) {
	repository, cleanup := setupPostgresTest(t)
	defer cleanup()

	// BIGSERIAL IDs are positive, so -1 cannot identify a stored task.
	missingID := -1

	_, err := repository.Get(missingID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Get(): expected ErrNotFound, got %v", err)
	}

	_, err = repository.Update(Task{
		ID:       missingID,
		Title:    testTaskTitle("missing"),
		Status:   "todo",
		Priority: "medium",
	})
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Update(): expected ErrNotFound, got %v", err)
	}

	err = repository.Delete(missingID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("Delete(): expected ErrNotFound, got %v", err)
	}
}
