FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Start a new stage from scratch
FROM alpine:latest

# Add non-root user
RUN adduser -D appuser

# Create app directory and set permissions
WORKDIR /app
COPY --from=builder /app/main .
COPY init.sql .

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Command to run the executable
CMD ["./main"] 