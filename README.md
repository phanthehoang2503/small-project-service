# Go Microservices E-Commerce (Practice Project)

This project is a **practice-only microservices backend** built with **Go**, **Gin**, **GORM**, and **Docker**.
It’s not for production use — it’s designed for learning how real systems are structured, communicate, and scale.

---

## Overview

Each service runs independently with its own routes, database models, and Swagger docs.
All services communicate through REST and RabbitMQ.

```mermaid
sequenceDiagram
    participant User
    participant Order as Order Service
    participant Product as Product Service
    participant Payment as Payment Service
    
    User->>Order: 1. Create Order (Pending)
    
    par Process Order
        Order-)Product: 2a. Event: "order.requested"
        Order-)Payment: 2b. Event: "order.requested"
    end
    
    alt Payment Success
        Payment-)Order: 3. Event: "order.paid"
        Order->>Order: 4. Update Status: Paid
    else Payment Failed
        Payment-)Order: 3. Event: "payment.failed"
        Order->>Order: 4. Cancel Order
    end
```
### Folder structure
```bash
D:.
├───auth-service
├───broker-service
├───cart-service
├───common
├───internal
├───logger-service
├───mailer-service
├───order-service
├───payment-service
└───product-service
```

### Key Features
*   **Event-Driven Architecture**: Uses RabbitMQ for asynchronous processing.
*   **Distributed Transactions (Saga Pattern)**: Ensures data consistency across services.
*   **Redis**: Accelerates product data retrieval.
*   **Resilient Messaging**: Broker automatically reconnects on network loss.
*   **Observability**: Centralized logging with Grafana Loki.
*   **CI/CD**: Automated build and test with GitHub Actions.

[Product service](product-service/README.md#product-service)  
[Cart service](cart-service/README.md#cart-service)  
[Order service](order-service/README.md#order-service)  
[Auth service](auth-service/README.md#auth-service)  
[Payment service](payment-service/README.md#payment-service)

---

### Current services

| Service | Port | Description |
|----------|------|-------------|
| **Product Service** | 8081 | Manages products and their CRUD operations |
| **Cart Service** | 8082 | Handles shopping cart creation, updates, and item management |
| **Order Service** | 8083 | Processes orders and connects with cart + user data |
| **Auth Service** | 8084 | Handles user registration, login, and JWT token generation |
| **Logger Service** | 8085 | Centralized logging service (gRPC/RabbitMQ consumer) |
| **Payment Service** | 8086 | Simulates payment processing |

### Infrastructure

| Service | Port | Description |
|----------|------|-------------|
| **PostgreSQL** | 5432 | Main database |
| **RabbitMQ** | 5672 | Message broker for async communication |
| **MailHog** | 1025 | Email testing tool |

---

### Tech stack

- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM (PostgreSQL)
- **Docs:** Swagger (`swaggo/swag`)
- **Containerization:** Docker & Docker Compose
- **Live Reload:** Air
- **Authentication:** JWT
- **Messaging:** RabbitMQ
- **Caching:** Redis
- **Observability:** Grafana, Loki, Promtail

---

### Setup

```bash
# Run all services
docker compose up --build

# Run individual service, example:
cd product-service
air
```

Then visit Swagger UI for each service:

***Product***: http://localhost:8081/swagger/index.html  
***Cart***: http://localhost:8082/swagger/index.html  
***Order***: http://localhost:8083/swagger/index.html  
***Auth***: http://localhost:8084/swagger/index.html  
***Payment***: http://localhost:8086/swagger/index.html

### How to use
1. Open the Auth page then register and login to get the token.
2. On cart or order page click on the lock icon and type: **Bearer <"token">** remove the (", <>) symbol.


## Author

**Thế Hoàng or you can call me *Josh*. Why? I just love that name and it shorter than my real name**  
CS Student | Backend Developer in training  
Learning Go, microservices, Java and JS.

# Vietnamese

## Xây dựng Microservices E-Commerce với Go

### Tổng quan

Mỗi service chạy độc lập, có route, model, và tài liệu Swagger riêng.
Tất cả giao tiếp với nhau qua REST và RabbitMQ.

### Về tính năng
*   **Kiến trúc hướng sự kiện**: Sử dụng RabbitMQ để xử lý bất đồng bộ.
*   **Giao dịch phân tán**: Đảm bảo tính nhất quán dữ liệu giữa các service.
*   **Cache**: Tăng tốc độ đọc dữ liệu sản phẩm.
*   **Resilient**: Broker tự động kết nối lại khi mất mạng, đảm bảo không mất tin nhắn.
*   **Giám sát (Observability)**: Ghi log tập trung với Grafana Loki.
*   **CI/CD**: Tự động build và test với GitHub Actions.

---

### Các service hiện có

| Service | Port | Mô tả |
|----------|------|-------|
| **Product Service** | 8081 | Quản lý sản phẩm |
| **Cart Service** | 8082 | Quản lý giỏ hàng |
| **Order Service** | 8083 | Xử lý đơn hàng |
| **Auth Service** | 8084 | Xử lý đăng ký, đăng nhập |
| **Logger Service** | 8085 | Ghi log |
| **Payment Service** | 8086 | Mô phỏng thanh toán |

### Infrastructure

| Service | Port | Mô tả |
|----------|------|-------|
| **PostgreSQL** | 5432 | CSDL |
| **RabbitMQ** | 5672 | Message |
| **MailHog** | 1025 | Dùng để test email |

---

### Công nghệ sử dụng

- **Ngôn ngữ:** Go
- **Framework:** Gin
- **ORM:** GORM (PostgreSQL)
- **Docs:** Swagger (`swaggo/swag`)
- **Container:** Docker & Docker Compose
- **Live Reload:** Air
- **Xác thực:** JWT
- **Messaging:** RabbitMQ
- **Cache:** Redis
- **Giám sát:** Grafana, Loki, Promtail

---

### Cách chạy

```bash
# Chạy tất cả service
docker compose up --build

# Chạy từng service riêng
cd auth-service
air
```

Truy cập Swagger UI của từng service:
***Product***: http://localhost:8081/swagger/index.html  
***Cart***: http://localhost:8082/swagger/index.html  
***Order***: http://localhost:8083/swagger/index.html  
***Auth***: http://localhost:8084/swagger/index.html  
***Payment***: http://localhost:8086/swagger/index.html  
[Cách sử dụng](#how-to-use)
