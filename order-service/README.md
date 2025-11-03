## order-service

This service handles order creation and order lifecycle for the small-project
microservices system.

### Overview

The `order-service` implements endpoints for creating orders, listing orders,
and retrieving order details. It follows the common project layout with an
entrypoint at `cmd/api/main.go` and internal packages for handlers, models,
repositories, and routing. Swagger docs are under `docs/`.

### Prerequisites

- Go 1.20+
- Docker & Docker Compose 

### Run locally

```powershell
cd order-service/cmd/api
go run .
```

### API Endpoints 

- POST /orders — create a new order
- GET /orders/{id} — get order by id
- GET /orders?userId={userId} — list user orders

### Swagger / API docs

http://localhost:8082/swagger/index.html#/Order/