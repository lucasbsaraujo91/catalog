
# Dockerfile-httpserver
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o httpserver ./cmd/httpserver

FROM debian:bookworm-slim
COPY --from=builder /app/httpserver /usr/local/bin/
EXPOSE 8080
ENTRYPOINT ["httpserver"]

