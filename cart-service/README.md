## cart-service

This service manages user shopping carts. It maintains a local copy of product data to ensure cart operations are fast and independent.

### Overview

The `cart-service` allows users to add, remove, and view items in their cart. It listens to product events to keep its local product data in sync with the `product-service`.

**Key Features:**
*   **Event-Driven Sync**: Updates local product table when `product-service` publishes changes.
*   **User Isolation**: Carts are linked to User IDs.

### Run locally

From repository root:

```powershell
cd cart-service/cmd/api
go run .
```

### API Endpoints

- GET /cart — get current user's cart
- POST /cart — add item to cart
- DELETE /cart/{id} — remove item from cart

### Events

- **Consumes**: `product.created`, `product.updated`, `product.deleted`

### Swagger / API docs

http://localhost:8082/swagger/index.html#/Cart/