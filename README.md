# Go Microservices E-Commerce (Practice Project)

This project is a **practice-only microservices backend** built with **Go**, **Gin**, **GORM**, and **Docker**.  
It’s not for production use — it’s designed for learning how real systems are structured, communicate, and scale.

---

## Overview

Each service runs independently with its own routes, database models, and Swagger docs.  
All services communicate through REST for simplicity (message queues can be added later).

```bash
D:.
├───api-test
│   ├───cart-api.http
│   ├───product-api
│   └───user-api
├───auth-service
│   ├───cmd
│   │   └───api
│   ├───docs
│   └───internal
│       ├───handler
│       ├───model
│       ├───repo
│       └───router
├───cart-service
│   ├───cmd
│   │   └───api
│   ├───docs
│   ├───internal
│   │   ├───handler
│   │   ├───model
│   │   ├───repo
│   │   └───router
│   └───tmp
├───common
│   ├───auth
│   └───logger
├───customer-service
├───internal
│   ├───database
│   ├───middleware
│   └───util
├───order-service
│   ├───.idea
│   ├───cmd
│   │   └───api
│   ├───docs
│   ├───internal
│   │   ├───handler
│   │   ├───model
│   │   ├───repo
│   │   └───router
│   └───tmp
└───product-service
    ├───cmd
    │   └───api
    ├───docs
    ├───internal
    │   ├───handler
    │   ├───model
    │   ├───repo
    │   └───router
    └───tmp
```
[Product service](product-service/README.md#product-service)  
[Cart service](cart-service/README.md#cart-service)  
[Order service](order-service/README.md#order-service)  
[Auth service](auth-service/README.md#auth-service)  

---

### Current services

| Service | Description |
|----------|-------------|
| **Auth Service** | Handles user registration, login, and JWT token generation |
| **Product Service** | Manages products and their CRUD operations |
| **Cart Service** | Handles shopping cart creation, updates, and item management |
| **Order Service** | Processes orders and connects with cart + user data |


### Learning goals

Practice **Go microservice structure** (using `internal/`, `cmd/`, `repo/`, `handler/`, etc.)
Use **Swagger** for API documentation
Manage multiple services via **Docker Compose**
Explore **auth, inter-service communication, and clean architecture**

---

### Tech stack

- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM (PostgreSQL)
- **Docs:** Swagger (`swaggo/swag`)
- **Containerization:** Docker & Docker Compose
- **Live Reload:** Air
- **Authentication:** JWT

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
Tất cả giao tiếp với nhau qua REST (sau này có thể thêm message queue).  
Cấu trúc file: [đầu trang](#overview)

---

### Các service hiện có

| Service | Mô tả |
|----------|-------|
| **Auth Service** | Xử lý đăng ký, đăng nhập, và tạo JWT token |
| **Product Service** | Quản lý sản phẩm (CRUD) |
| **Cart Service** | Quản lý giỏ hàng của người dùng |
| **Order Service** | Xử lý đơn hàng, kết nối giỏ hàng và người dùng |

---

### Dự kiến sẽ thêm vào

| Service | Mô tả |
|----------|----------------------|
| **Inventory Service** | Quản lý tồn kho |
| **Payment Service** | Mô phỏng thanh toán |
| **Notification Service** | Gửi thông báo(email/SMS)|
| **API Gateway** | Tạo điểm truy cập và xác thực|

---

### Các kiến thức học được

Thực hành **cấu trúc microservice trong Go** (`internal/`, `cmd/`, `repo/`, `handler/`, v.v.)  
Dùng **Swagger** để tạo tài liệu API  
Quản lý nhiều service với **Docker Compose**  
Tìm hiểu **xác thực, giao tiếp giữa các service, và clean architecture**  

---

## Công nghệ sử dụng

- **Ngôn ngữ:** Go  
- **Framework:** Gin  
- **ORM:** GORM (PostgreSQL)  
- **Docs:** Swagger (`swaggo/swag`)  
- **Container:** Docker & Docker Compose  
- **Live Reload:** Air  
- **Xác thực:** JWT  

---

## Cách chạy

```bash
# Chạy tất cả service
docker compose up --build

# Chạy từng service riêng lẻ
cd auth-service
air
```

Truy cập Swagger UI của từng service:  
***Product***: http://localhost:8081/swagger/index.html  
***Cart***: http://localhost:8082/swagger/index.html  
***Order***: http://localhost:8083/swagger/index.html  
***Auth***: http://localhost:8084/swagger/index.html  
[Cách sử dụng](#how-to-use)

