# Build stage
FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o api ./cmd/api
RUN go build -o worker ./cmd/worker
RUN go build -o simulator ./cmd/simulator

# API Run stage
FROM alpine:latest AS api
WORKDIR /app
COPY --from=builder /app/api .
EXPOSE 8000
CMD ["./api"]

# Worker Run stage
FROM alpine:latest AS worker
WORKDIR /app
COPY --from=builder /app/worker .
CMD ["./worker"]

# Simulator Run stage
FROM alpine:latest AS simulator
WORKDIR /app
COPY --from=builder /app/simulator .
CMD ["./simulator"]
