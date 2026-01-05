# base go image
FROM golang:1.25.3-alpine3.22 AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
# Build the binary pointing to the package path
RUN go build -o /app/gateway-service ./gateway-service/main.go

# tiny image to runs the binary
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/gateway-service /app/gateway-service

EXPOSE 8080
CMD ["/app/gateway-service"]