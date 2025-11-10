# base go image
FROM golang:1.25.3-alpine3.22 AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
ENV CGO_ENABLED=0
RUN go build -o /app/logger-service ./logger-service/cmd/api

#tiny image to runs the binary
FROM alpine:3.22
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/logger-service /app/logger-service

EXPOSE 8085
CMD ["/app/logger-service"]