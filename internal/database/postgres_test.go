package database

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"
)

func TestNewPostgresPoolSuccess(t *testing.T) {
	databaseURL := os.Getenv("TEST_DATABASE_URL")
	if databaseURL == "" {
		t.Skip("TEST_DATABASE_URL is not set")
	}

	// Prevent accidentally connecting to the development database.
	if !strings.Contains(databaseURL, "/task_manager_test?") {
		t.Fatal("TEST_DATABASE_URL must point to task_manager_test")
	}

	pool, err := NewPostgresPool(databaseURL)
	if err != nil {
		t.Fatalf("expected successful connection, got: %v", err)
	}
	defer pool.Close()

	if pool == nil {
		t.Fatal("expected a non-nil connection pool")
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	var result int

	err = pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		t.Fatalf("SELECT 1 failed: %v", err)
	}

	if result != 1 {
		t.Fatalf("expected 1, got %d", result)
	}
}

func TestNewPostgresPoolInvalidURL(t *testing.T) {
	pool, err := NewPostgresPool(
		"postgres://user:password@localhost:invalid/database",
	)

	if err == nil {
		if pool != nil {
			pool.Close()
		}
		t.Fatal("expected configuration error, got nil")
	}

	if pool != nil {
		pool.Close()
		t.Fatal("expected nil pool for invalid URL")
	}

	if !strings.Contains(err.Error(), "parse database configuration") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNewPostgresPoolConnectionFailure(t *testing.T) {
	// Port 1 on localhost should normally have no PostgreSQL server.
	// The test expects the connection to fail.
	databaseURL := "postgres://invalid:invalid@127.0.0.1:1/invalid?sslmode=disable"

	start := time.Now()

	pool, err := NewPostgresPool(databaseURL)

	elapsed := time.Since(start)

	if pool != nil {
		pool.Close()
		t.Fatal("expected nil pool for failed connection")
	}

	if err == nil {
		t.Fatal("expected connection error, got nil")
	}

	if !strings.Contains(err.Error(), "connect to PostgreSQL") {
		t.Fatalf("unexpected error: %v", err)
	}

	if elapsed > 7*time.Second {
		t.Fatalf("connection took too long: %v", elapsed)
	}
}
