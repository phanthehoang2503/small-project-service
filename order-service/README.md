## order-service

This service handles order placement, retrieval, and status management. It is the central coordinator for the **Saga Pattern**.

### Overview

The `order-service` is a Go HTTP API that manages the lifecycle of orders. It communicates with `product-service` and `payment-service` via RabbitMQ to ensure data consistency.

### Run locally

From repository root:

```powershell
cd order-service/cmd/api
go run .
```

Ensure RabbitMQ and Postgres are running.

### API Endpoints

- GET /orders — list orders (User specific)
- GET /orders/{id} — get order by UUID
- GET /orders/search?id={id} — get order by numeric ID
- POST /orders — create order (Triggers `order.requested` event)

### Events

- **Publishes**: `order.requested`
- **Consumes**: `order.paid`, `payment.failed`, `stock.failed`

### Swagger / API docs

http://localhost:8083/swagger/index.html#/Orders/