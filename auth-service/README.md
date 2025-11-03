## auth-service

This service handles authentication and user management for the small-project
microservices system.

### Overview

The `auth-service` provides endpoints for user registration, login, token
issuance, and basic user CRUD operations. Entry point is `cmd/api/main.go` and
the service uses `internal/handler`, `internal/model`, `internal/repo`, and
`internal/router` packages. Swagger docs are available under `docs/`.

### Prerequisites

- Go 1.20+ installed
- Docker & Docker Compose

### Run locally

```powershell
cd auth-service/cmd/api
go run .
```

### API Endpoints

- POST /auth/register — register a new user
- POST /auth/login — obtain JWT
- GET /users/{id} — get user details

### Swagger / API docs

http://localhost:8084/swagger/index.html#/Auth/