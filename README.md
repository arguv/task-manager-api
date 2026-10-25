# Task Manager API

A RESTful task management API built with Go's standard `net/http` package.

## Project Overview

A backend application for creating, reading, updating, and deleting tasks.

This project is being developed as a practical learning experience, following professional Go backend engineering practices.

## Technology Stack

- Go 1.27.3+
- `net/http`
- In-memory task storage

## Current Features

- Task CRUD operations
- Task status and priority management
- JSON request and response handling
- Health check endpoint

Valid task statuses are `todo`, `in_progress`, and `done`. Valid priorities are
`low`, `medium`, and `high`.

## API Endpoints

| Method | Endpoint | Description |
| --- | --- | --- |
| `GET` | `/health` | Returns the service health status |
| `GET` | `/api/v1/tasks` | Lists all tasks |
| `GET` | `/api/v1/tasks/{id}` | Returns one task |
| `POST` | `/api/v1/tasks` | Creates a task |
| `PUT` | `/api/v1/tasks/{id}` | Updates a task |
| `DELETE` | `/api/v1/tasks/{id}` | Deletes a task |

Task creation and update requests use this JSON shape:

```json
{
  "title": "Review API documentation",
  "status": "todo",
  "priority": "medium"
}
```

Example request:

```bash
curl http://localhost:8080/api/v1/tasks
```

## Project Structure

```text
task-manager-api/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
├── migrations/
├── tests/
├── .gitignore
├── go.mod
└── README.md
```

## Development Status

In development — the current implementation stores tasks in memory, so data is
reset whenever the server restarts.

## Getting Started

The module path is `github.com/arguv/task-manager-api`. Make sure Go 1.27.3 or
newer is installed, then format and validate the source code:

```bash
gofmt -w cmd/api/main.go
go vet ./...
```

Download module dependencies, if any are added:

```bash
go mod download
```

Run the API with:

```bash
go run ./cmd/api
```

The server starts on port `8080`.

To run the test suite:

```bash
go test ./...
```

## License

No license has been selected yet.
