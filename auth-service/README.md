## auth-service

This service handles user authentication and authorization.

### Overview

The `auth-service` manages user registration and login. It issues **JWT (JSON Web Tokens)** that are used to authenticate requests to other services.

**Key Features:**
*   **JWT Issuance**: Secure stateless authentication.
*   **Password Hashing**: Uses bcrypt for security.

### Run locally

From repository root:

```powershell
cd auth-service/cmd/api
go run .
```

### API Endpoints

- POST /auth/register — register new user
- POST /auth/login — login and get JWT token

### Swagger / API docs

http://localhost:8084/swagger/index.html#/Auth/