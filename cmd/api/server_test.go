package main

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewHTTPServer(t *testing.T) {
	handler := http.NewServeMux()

	server := newHTTPServer(":8080", handler)

	if server.Addr != ":8080" {
		t.Fatalf("expected :8080, got %s", server.Addr)
	}

	if server.Handler != handler {
		t.Fatal("expected the provided HTTP handler")
	}

	if server.ReadHeaderTimeout != 5*time.Second {
		t.Fatalf(
			"expected 5s ReadHeaderTimeout, got %v",
			server.ReadHeaderTimeout,
		)
	}

	if server.IdleTimeout != 60*time.Second {
		t.Fatalf(
			"expected 60s IdleTimeout, got %v",
			server.IdleTimeout,
		)
	}
}

func TestHealthCheck(t *testing.T) {
	request := httptest.NewRequest(
		http.MethodGet,
		"/health",
		nil,
	)

	recorder := httptest.NewRecorder()

	healthCheck(recorder, request)

	response := recorder.Result()
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf(
			"expected HTTP 200, got %d",
			response.StatusCode,
		)
	}
}

func TestHTTPServerGracefulShutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	server := newHTTPServer(":0", handler)

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("create listener: %v", err)
	}

	serveDone := make(chan error, 1)

	go func() {
		serveDone <- server.Serve(listener)
	}()

	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	response, err := client.Get("http://" + listener.Addr().String())
	if err != nil {
		t.Fatalf("send request: %v", err)
	}
	response.Body.Close()

	if response.StatusCode != http.StatusOK {
		t.Fatalf("expected HTTP 200, got %d", response.StatusCode)
	}

	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		t.Fatalf("graceful shutdown failed: %v", err)
	}

	select {
	case err := <-serveDone:
		if err != http.ErrServerClosed {
			t.Fatalf("unexpected server error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("server did not stop after shutdown")
	}
}
