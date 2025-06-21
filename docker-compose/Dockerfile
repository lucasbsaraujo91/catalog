# Build stage
FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN go build -o ordersystem ./cmd/ordersystem

# Final stage
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /app/ordersystem /usr/local/bin/

EXPOSE 8080

ENTRYPOINT ["ordersystem"]

