# Build stage
FROM golang:1.22.2-alpine AS builder

# Install git for dependency fetching
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files from root (build context is now root)
COPY services/product-service/go.mod services/product-service/go.sum ./services/product-service/
COPY shared/ ./shared/

# Set working directory to service
WORKDIR /app/services/product-service

# Download dependencies
RUN go mod download

# Copy source code
COPY services/product-service/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main cmd/main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/product-service/main .

# Expose port 8082
EXPOSE 8082

# Command to run
CMD ["./main"]
