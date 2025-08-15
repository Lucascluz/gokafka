#!/bin/bash

# Master script to stop the entire gokafka system

PROJECT_ROOT="/home/lucas/Projects/gokafka"
SCRIPTS_DIR="$PROJECT_ROOT/scripts"

echo "🛑 Stopping GoKafka System..."
echo "============================="

# Stop Go services first
echo "🔧 Stopping Go services..."
bash "$SCRIPTS_DIR/stop-services.sh"

# Stop infrastructure
echo "📦 Stopping infrastructure services..."
bash "$SCRIPTS_DIR/stop-infrastructure.sh"

echo ""
echo "🎯 GoKafka system stopped successfully!"
