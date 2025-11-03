## cart-service

This service manages shopping cart operations for the small-project system.

### Overview

The `cart-service` provides endpoints to create and manage user carts and cart
items. Entrypoint: `cmd/api/main.go`. Internals follow the same layout used by
other services: `internal/handler`, `internal/model`, `internal/repo`, and
`internal/router`. Swagger docs live in `docs/`.

### Prerequisites

- Go 1.20+
- Docker & Docker Compose

### Run locally

```powershell
cd cart-service/cmd/api
go run .
```

### API Endpoints

- POST /carts — create a cart
- GET /carts/{userId} — get user's cart
- POST /carts/{userId}/items — add item
- DELETE /carts/{userId}/items/{itemId} — remove item


### Swagger / API docs

http://localhost:8082/swagger/index.html#/Cart