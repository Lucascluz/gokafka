# Build stage
FROM golang:1.22.2-alpine AS builder

# Install git for dependency fetching
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files from root (build context is now root)
COPY services/user-service/go.mod services/user-service/go.sum ./services/user-service/
COPY shared/ ./shared/

# Set working directory to service
WORKDIR /app/services/user-service

# Download dependencies
RUN go mod download

# Copy source code
COPY services/user-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/user-service/main .

# Command to run
CMD ["./main"]
