# base go image
FROM golang:1.25.3-alpine3.22 AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/product-service ./product-service/cmd/api

#tiny image to runs the binary
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/product-service /app/product-service

EXPOSE 8081
CMD ["/app/product-service"]