# Build stage
FROM golang:1.22.2-alpine AS builder

# Install git for dependency fetching
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files from root (build context is now root)
COPY services/api-gateway/go.mod services/api-gateway/go.sum ./services/api-gateway/
COPY shared/ ./shared/

# Set working directory to service
WORKDIR /app/services/api-gateway

# Download dependencies
RUN go mod download

# Copy source code
COPY services/api-gateway/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/api-gateway/main .

# Expose port 8080
EXPOSE 8080

# Command to run
CMD ["./main"]
