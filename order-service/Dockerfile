# base go image
FROM golang:1.25.3-alpine3.22 AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY internal ./internal
COPY order-service ./order-service

ENV CGO_ENABLED=0
RUN go build -o /app/order-service ./order-service/cmd/api

#tiny image to runs the binary
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /app/order-service /app/order-service

EXPOSE 8083
CMD ["/app/order-service"]