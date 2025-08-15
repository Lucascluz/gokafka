#!/bin/bash

# Script to start all Go services

PROJECT_ROOT="/home/lucas/Projects/gokafka"

echo "Starting all Go services..."

# Function to start a service in the background
start_service() {
    local service_name=$1
    local service_path=$2
    
    echo "Starting $service_name..."
    cd "$PROJECT_ROOT/$service_path"
    
    # Kill any existing process on the same port (optional)
    case $service_name in
        "api-gateway")
            lsof -ti:8080 | xargs kill -9 2>/dev/null || true
            ;;
        "user-service")
            # User service doesn't have an HTTP port, but we can kill any existing process
            pkill -f "user-service/cmd/main.go" 2>/dev/null || true
            ;;
    esac
    
    # Start the service in the background
    nohup go run cmd/main.go > "$PROJECT_ROOT/logs/${service_name}.log" 2>&1 &
    echo $! > "$PROJECT_ROOT/logs/${service_name}.pid"
    
    echo "$service_name started (PID: $(cat $PROJECT_ROOT/logs/${service_name}.pid))"
}

# Create logs directory if it doesn't exist
mkdir -p "$PROJECT_ROOT/logs"

# Start services
start_service "user-service" "services/user-service"
sleep 3  # Give user-service time to start

start_service "api-gateway" "services/api-gateway"
sleep 2  # Give api-gateway time to start

echo ""
echo "All services started successfully!"
echo ""
echo "Service status:"
echo "User Service - PID: $(cat $PROJECT_ROOT/logs/user-service.pid 2>/dev/null || echo 'Not found')"
echo "API Gateway - PID: $(cat $PROJECT_ROOT/logs/api-gateway.pid 2>/dev/null || echo 'Not found') - http://localhost:8080"
echo ""
echo "Logs are available in the logs/ directory:"
echo "- User Service: logs/user-service.log"
echo "- API Gateway: logs/api-gateway.log"
echo ""
echo "To stop all services, run: ./scripts/stop-services.sh"
