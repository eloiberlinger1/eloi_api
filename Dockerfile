# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy dependency manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and config files
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -o myapp ./cmd

# Execution stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and data.json from the builder
COPY --from=builder /app/myapp .
COPY --from=builder /app/data.json .

# Expose port (Coolify can map this)
EXPOSE 8080

# Run the app
CMD ["./myapp"]
