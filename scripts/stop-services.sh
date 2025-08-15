#!/bin/bash

# Script to stop all Go services

PROJECT_ROOT="/home/lucas/Projects/gokafka"

echo "Stopping all Go services..."

# Function to stop a service
stop_service() {
    local service_name=$1
    local pid_file="$PROJECT_ROOT/logs/${service_name}.pid"
    
    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            echo "Stopping $service_name (PID: $pid)..."
            kill $pid
            sleep 2
            
            # Force kill if still running
            if ps -p $pid > /dev/null 2>&1; then
                echo "Force killing $service_name..."
                kill -9 $pid
            fi
        else
            echo "$service_name was not running"
        fi
        rm -f "$pid_file"
    else
        echo "No PID file found for $service_name"
    fi
}

# Stop services
stop_service "api-gateway"
stop_service "user-service"

# Also kill any remaining processes by name (backup method)
echo "Cleaning up any remaining processes..."
pkill -f "api-gateway/cmd/main.go" 2>/dev/null || true
pkill -f "user-service/cmd/main.go" 2>/dev/null || true

# Kill processes on specific ports
lsof -ti:8080 | xargs kill -9 2>/dev/null || true

echo "All services stopped!"
