## payment-service

This service handles payment processing (mocked) for orders.

### Overview

The `payment-service` is a background worker that listens for order requests. It simulates a payment gateway interaction and updates the order status.

**Key Features:**
*   **Mock Payment Gateway**: Simulates success/failure scenarios.
*   **Saga Participant**: Publishes success/failure events to drive the order workflow.

### Run locally

From repository root:

```powershell
cd payment-service/cmd/api
go run .
```

### Events

- **Consumes**: `order.requested`
- **Publishes**: `order.paid`, `payment.failed`

### Swagger / API docs

http://localhost:8086/swagger/index.html
