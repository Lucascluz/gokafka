#!/bin/bash

# Master script to start the entire gokafka system

PROJECT_ROOT="/home/lucas/Projects/gokafka"
SCRIPTS_DIR="$PROJECT_ROOT/scripts"

echo "🚀 Starting GoKafka System..."
echo "=========================="

# Make sure all scripts are executable
chmod +x "$SCRIPTS_DIR"/*.sh

# Start infrastructure
echo "📦 Starting infrastructure services..."
bash "$SCRIPTS_DIR/start-infrastructure.sh"

if [ $? -eq 0 ]; then
    echo "✅ Infrastructure services started successfully!"
    
    echo ""
    echo "⏳ Waiting 20 seconds for infrastructure to be fully ready..."
    sleep 20
    
    # Start Go services
    echo "🔧 Starting Go services..."
    bash "$SCRIPTS_DIR/start-services.sh"
    
    if [ $? -eq 0 ]; then
        echo ""
        echo "🎉 GoKafka system is now running!"
        echo "================================"
        echo ""
        echo "🌐 API Gateway: http://localhost:8080"
        echo "🧪 Test endpoint: http://localhost:8080/test"
        echo ""
        echo "📊 Infrastructure:"
        echo "   PostgreSQL: localhost:5432"
        echo "   Redis: localhost:6379"
        echo "   Kafka: localhost:9092"
        echo ""
        echo "📝 Logs are available in the logs/ directory"
        echo ""
        echo "🛑 To stop everything: ./scripts/stop-all.sh"
    else
        echo "❌ Failed to start Go services"
        exit 1
    fi
else
    echo "❌ Failed to start infrastructure services"
    exit 1
fi
