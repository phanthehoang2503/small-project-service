## product-service

This service manages products (CRUD) for the small-project microservices system.

### Overview

### Overview

The `product-service` is a Go HTTP API that provides endpoints to create, list,
retrieve, and delete products. It now includes:
*   **Redis Caching**: To improve performance for read operations.
*   **RabbitMQ Consumer**: To handle async stock deduction when orders are placed.

The service follows the repository layout used across the project (entrypoint at
`cmd/api/main.go`, internal packages for handlers, models, repo, and the router).
Swagger docs are included in the `docs/` folder.

### Run locally

From repository root you can run the service directly with Go. Open PowerShell
and run:

```powershell
cd product-service/cmd/api
go run .
```

By default the service will read configuration from environment variables.
Check `cmd/api/main.go` and `internal/router/router.go` for exact port/env names.
Ensure you have **Redis** running locally or in Docker for caching to work.

### API Endpoints

The API endpoints mirror the `.http` files in `product-api/` (if present). Typical
endpoints are:

- GET /products — list products
- GET /products/{id} — get product by id (**Cached**)
- POST /products — create product (JSON body)
- PUT /products/{id} — update product (**Invalidates Cache**)
- DELETE /products/{id} — delete product (**Invalidates Cache**)

Example `curl` requests:

```powershell
# list
curl http://localhost:8081/products

# get by id (Checks Redis first)
curl http://localhost:8081/products/1

# create
curl -X POST http://localhost:8081/products -H "Content-Type: application/json" -d '{"name":"T-shirt","price":30000}'
```

### Swagger / API docs

http://localhost:8081/swagger/index.html#/Products/
