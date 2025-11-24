## payment-service

This service handles payment processing for the small-project microservices system.

### Overview

The `payment-service` provides endpoints to process payments for orders. It integrates with the `order-service` to update order status upon successful payment. Entrypoint: `cmd/api/main.go`. Internals follow the standard project layout: `internal/handler`, `internal/model`, `internal/repo`, and `internal/router`. Swagger docs are available in `docs/`.

### Prerequisites

- Go 1.20+
- Docker & Docker Compose

### Run locally

```powershell
cd payment-service/cmd/api
go run .
```

### API Endpoints

- POST /payments — process a payment for an order

### Swagger / API docs

http://localhost:8085/swagger/index.html#/Payment
