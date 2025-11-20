FROM golang:1.25.3-alpine3.22 AS builder
RUN apk add --no-cache git ca-certificates

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=0
RUN go build -o /app/payment-service ./payment-service/cmd/api


FROM alpine:3.22
RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /app/payment-service /app/payment-service

EXPOSE 8086

CMD ["/app/payment-service"]
