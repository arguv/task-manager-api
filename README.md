# Task Manager API

A RESTful task management API built with Go's standard `net/http` package.
The API currently uses an in-memory repository, so all data is reset when the
server restarts.

## Technology stack

- Go 1.27.3+
- `net/http` with method-aware `http.ServeMux` routing
- In-memory storage protected by a mutex
- JSON request and response bodies

## Features

- Create, list, retrieve, update, and delete tasks
- Task validation for title, status, and priority
- Health check endpoint
- Unit and HTTP handler/router tests

Valid statuses are `todo`, `in_progress`, and `done`. Valid priorities are
`low`, `medium`, and `high`.

## API endpoints

| Method | Endpoint | Success | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | `200` | Returns service health |
| `GET` | `/api/v1/tasks` | `200` | Lists all tasks |
| `GET` | `/api/v1/tasks/{id}` | `200` | Returns one task |
| `POST` | `/api/v1/tasks` | `201` | Creates a task |
| `PUT` | `/api/v1/tasks/{id}` | `200` | Replaces a task |
| `DELETE` | `/api/v1/tasks/{id}` | `204` | Deletes a task |

Create and update requests use this JSON shape:

```json
{
  "title": "Review API documentation",
  "status": "todo",
  "priority": "medium"
}
```

Invalid JSON, invalid task data, invalid IDs, and missing tasks return `4xx`
responses as appropriate.

## Project structure

```text
task-manager-api/
├── cmd/
│   └── api/
│       └── main.go             # Application entry point and health route
├── internal/
│   └── task/
│       ├── model.go            # Task request/response models
│       ├── repository.go       # Repository interface and memory storage
│       ├── service.go          # Validation and business logic
│       ├── handler.go          # HTTP handlers
│       ├── router.go           # API route registration
│       └── *_test.go           # Service, handler, router, and validation tests
├── .gitignore
├── go.mod
└── README.md
```

## Development commands

Run these commands from the repository root.

Check the installed Go version:

```bash
go version
```

Download module dependencies:

```bash
go mod download
```

Format the source code:

```bash
gofmt -w cmd/api/main.go internal/task/*.go
```

Run static analysis:

```bash
go vet ./...
```

Run all tests:

```bash
go test ./...
```

Generate a coverage profile:

```bash
go test -coverprofile=coverage.out ./...
```

Create an HTML coverage report:

```bash
go tool cover -html=coverage.out -o coverage.html
```

Run the API locally:

```bash
go run ./cmd/api
```

The server listens on `http://localhost:8080`. Stop it with `Ctrl+C`.

## Example requests

With the server running, use the following commands in another terminal:

```bash
# Health check
curl http://localhost:8080/health

# List tasks
curl http://localhost:8080/api/v1/tasks

# Get task 1
curl http://localhost:8080/api/v1/tasks/1

# Create a task
curl -X POST http://localhost:8080/api/v1/tasks \
  -H 'Content-Type: application/json' \
  -d '{"title":"Review API documentation","status":"todo","priority":"medium"}'

# Update task 1
curl -X PUT http://localhost:8080/api/v1/tasks/1 \
  -H 'Content-Type: application/json' \
  -d '{"title":"Learn advanced Go","status":"done","priority":"high"}'

# Delete task 1
curl -i -X DELETE http://localhost:8080/api/v1/tasks/1
```

## Development status

The API is in development. Storage is intentionally in memory for now; a
future persistence layer can implement the existing `Repository` interface.

## License

No license has been selected yet.
