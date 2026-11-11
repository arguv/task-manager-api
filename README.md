# Task Manager API

A RESTful task management API built with Go's standard `net/http` package and
PostgreSQL persistence.

## Technology stack

- Go 1.27.3+
- `net/http` with method-aware `http.ServeMux` routing
- PostgreSQL 17
- `pgx/v5` connection pooling
- JSON request and response bodies
- Docker Compose for local PostgreSQL development

## Features

- Create, list, retrieve, update, and delete tasks
- PostgreSQL-backed persistence
- Task validation for title, status, and priority
- Health check endpoint
- Graceful shutdown on `SIGINT` and `SIGTERM`
- Unit and HTTP handler/router tests
- Consistent JSON errors, panic recovery, request logging, optional API-key protection
- Pagination, filtering, sorting, OpenAPI documentation, Docker and GitHub Actions CI

Valid statuses are `todo`, `in_progress`, and `done`. Valid priorities are
`low`, `medium`, and `high`.

## API endpoints

| Method | Endpoint | Success | Description |
| --- | --- | --- | --- |
| `GET` | `/health` | `200` | Returns service health |
| `GET` | `/api/v1/tasks` | `200` | Lists tasks; supports `page`, `page_size`, `status`, `priority`, and `sort` |
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
responses as appropriate. Server failures return a JSON object such as
`{"error":"Internal server error"}` without exposing database details.

When `API_KEY` is configured, POST, PUT, and DELETE requests require the
`X-API-Key` header. GET requests remain public.

## Local setup

### Prerequisites

- Go 1.27.3 or newer
- Docker and Docker Compose
- `psql` for applying migrations

### Start PostgreSQL

Copy the example environment file and adjust the values if needed:

```bash
cp .env.example .env
```

The application also requires `DATABASE_URL`. Add it to `.env` for the default
local setup:

```dotenv
DATABASE_URL=postgres://task_user:change_me@localhost:5433/task_manager?sslmode=disable
```

Start the database container:

```bash
docker compose up -d postgres
```

Apply the initial schema:

```bash
PGPASSWORD=change_me psql \
  "postgres://task_user:change_me@localhost:5433/task_manager?sslmode=disable" \
  -f migrations/000001_create_tasks_table.up.sql
```

Run the API:

```bash
go run ./cmd/api
```

The server listens on `http://localhost:8080`. Stop it with `Ctrl+C`. Stop the
database container when finished with:

```bash
docker compose down
```

### Run the API with Docker

After starting PostgreSQL, build and run the API with:

```bash
docker build -t task-manager-api .
docker run --rm --network host --env-file .env task-manager-api
```

The complete endpoint contract is in [openapi.yaml](./openapi.yaml).

Use `docker compose down -v` only when you also want to remove the PostgreSQL
data volume.

## Development commands

Run these commands from the repository root.

```bash
go version
go mod download
gofmt -w cmd internal
go vet ./...
go test ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

The tests use in-memory and mock repositories where appropriate, so the full
test suite does not require a running PostgreSQL instance.

PostgreSQL integration tests run when `TEST_DATABASE_URL` points to a dedicated
database named `task_manager_test`; otherwise those tests are intentionally
skipped. CI can provide that variable after applying the migration.

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

## Project structure

```text
task-manager-api/
├── cmd/api/
│   ├── main.go               # Application entry point and health route
│   ├── server.go             # HTTP server configuration
│   └── server_test.go
├── internal/
│   ├── database/
│   │   └── postgres.go       # PostgreSQL connection pool
│   └── task/
│       ├── model.go          # Task request/response models
│       ├── repository.go     # Repository interface and memory repository
│       ├── postgres_repository.go
│       ├── service.go        # Validation and business logic
│       ├── handler.go        # HTTP handlers
│       ├── router.go         # API route registration
│       └── *_test.go         # Unit and HTTP tests
├── migrations/               # SQL schema migrations
├── openapi.yaml              # OpenAPI 3 endpoint contract
├── Dockerfile                # Multi-stage production image
├── compose.yaml              # Local PostgreSQL service
├── .env.example              # Local database configuration template
├── go.mod
└── README.md
```

## Development status

The API is in development. PostgreSQL is the runtime persistence layer, while
the repository interface and in-memory implementation remain useful for unit
tests and alternate storage implementations.

## License

No license has been selected yet.
