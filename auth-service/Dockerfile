# build stage
FROM golang:1.25.3-alpine3.22 AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/auth-service ./auth-service/cmd/api

# runtime stage
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/auth-service /app/auth-service

EXPOSE 8084
CMD ["/app/auth-service"]