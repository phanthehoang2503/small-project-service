# Go Microservices E-Commerce (Practice Project)

This project is a **practice-only microservices backend** built with **Go**, **Gin**, **GORM**, and **Docker**.
It’s not for production use — it’s designed for learning how real systems are structured, communicate, and scale.

---

## Overview

Each service runs independently with its own routes, database models, and Swagger docs.
All services communicate through REST and RabbitMQ.

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
| **RabbitMQ** | 5672 (15672 UI) | Message broker for async communication |
| **MailHog** | 1025 (8025 UI) | Email testing tool |

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
| **PostgreSQL** | 5432 | Database chính |
| **RabbitMQ** | 5672 (15672 UI) | Message broker |
| **MailHog** | 1025 (8025 UI) | Công cụ test email |

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
