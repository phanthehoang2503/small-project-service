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

---

## Current Services

| Service | Description |
|----------|-------------|
| **Auth Service** | Handles user registration, login, and JWT token generation |
| **Product Service** | Manages products and their CRUD operations |
| **Cart Service** | Handles shopping cart creation, updates, and item management |
| **Order Service** | Processes orders and connects with cart + user data |

---

## Planned Additions

| Planned Service | What You'll Learn |
|------------------|------------------|
| **Inventory Service** | Stock tracking, concurrency, and atomic updates |
| **Payment Service** | Simulated payment flow & external API integration |
| **Notification Service** | Async message processing (email/SMS mock) |
| **API Gateway** | Unified entrypoint & centralized auth validation |
| **Review Service** | Cross-service data relations (User + Product) |

---

## Learning Goals

- Practice **Go microservice structure** (using `internal/`, `cmd/`, `repo/`, `handler/`, etc.)
- Use **Swagger** for API documentation
- Manage multiple services via **Docker Compose**
- Explore **auth, inter-service communication, and clean architecture**
- Eventually experiment with **async messaging** (RabbitMQ/Kafka) and **caching** (Redis)

---

## Tech Stack

- **Language:** Go
- **Framework:** Gin
- **ORM:** GORM (PostgreSQL)
- **Docs:** Swagger (`swaggo/swag`)
- **Containerization:** Docker & Docker Compose
- **Live Reload:** Air
- **Authentication:** JWT

---

## Setup

```bash
# Run all services
docker compose up --build

# Run individual service, example:
cd product-service
air
```
Then visit Swagger UI for each service:
- Product: http://localhost:8081/swagger/index.html
- Cart: http://localhost:8082/swagger/index.html
- Order: http://localhost:8083/swagger/index.html
- Auth: http://localhost:8084/swagger/index.html

---

## Note

> This repository is for **learning purpose only**.  
> The goal is to understand architecture and clean service design, not production-grade performance or security.

---

## Author

**Hoang (phanthehoang2503)**  
- CS Student | Backend Developer in training  
- Learning Go, microservices, and AI automation

# Vietnamese

# Dự án Thực Hành Microservices E-Commerce bằng Go

Dự án này là **backend microservices chỉ dùng cho mục đích học tập**, được xây dựng bằng **Go**, **Gin**, **GORM**, và **Docker**.  
Không dành cho môi trường production — mục tiêu là để học cách các hệ thống thực tế được tổ chức, giao tiếp và mở rộng.

---

## Tổng quan

Mỗi service chạy độc lập, có route, model, và tài liệu Swagger riêng.  
Tất cả giao tiếp với nhau qua REST (sau này có thể thêm message queue).

---

## Các Service Hiện Có

| Service | Mô tả |
|----------|-------|
| **Auth Service** | Xử lý đăng ký, đăng nhập, và tạo JWT token |
| **Product Service** | Quản lý sản phẩm (CRUD) |
| **Cart Service** | Quản lý giỏ hàng của người dùng |
| **Order Service** | Xử lý đơn hàng, kết nối giỏ hàng và người dùng |

---

## Dự Kiến Thêm Trong Tương Lai

| Service | Kiến thức sẽ học được |
|----------|----------------------|
| **Inventory Service** | Quản lý tồn kho, xử lý song song và cập nhật nguyên tử |
| **Payment Service** | Mô phỏng thanh toán, tích hợp API ngoài |
| **Notification Service** | Gửi thông báo bất đồng bộ (email/SMS giả lập) |
| **API Gateway** | Tạo điểm truy cập chung và xác thực tập trung |
| **Review Service** | Kết nối dữ liệu giữa các service (User + Product) |

---

## Mục Tiêu Học Tập

- Thực hành **cấu trúc microservice trong Go** (`internal/`, `cmd/`, `repo/`, `handler/`, v.v.)  
- Dùng **Swagger** để tạo tài liệu API  
- Quản lý nhiều service với **Docker Compose**  
- Tìm hiểu **xác thực, giao tiếp giữa các service, và clean architecture**  
- Sau này thử nghiệm với **message queue (RabbitMQ/Kafka)** và **cache (Redis)**  

---

## Công Nghệ Sử Dụng

- **Ngôn ngữ:** Go  
- **Framework:** Gin  
- **ORM:** GORM (PostgreSQL)  
- **Docs:** Swagger (`swaggo/swag`)  
- **Container:** Docker & Docker Compose  
- **Live Reload:** Air  
- **Xác thực:** JWT  

---

## Cách Chạy

```bash
# Chạy tất cả service
docker compose up --build

# Chạy 1 service riêng
cd auth-service
air
```

Truy cập Swagger UI của từng service:
- Auth: http://localhost:8081/swagger/index.html  
- Product: http://localhost:8082/swagger/index.html  
- Cart: http://localhost:8083/swagger/index.html  
- Order: http://localhost:8084/swagger/index.html  

---

## Ghi chú

> Dự án này chỉ phục vụ **mục đích học tập**.  
> Mục tiêu là hiểu về kiến trúc và thiết kế service sạch.

---

## Tác giả

**Hoàng (phanthehoang2503)**  
💻 Sinh viên KHMT | Đang học backend development  
🚀 Học Go, microservices và tự động hóa AI
