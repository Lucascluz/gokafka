#!/bin/bash

# Script to stop all infrastructure services

echo "Stopping infrastructure services..."

# Stop and remove containers
echo "Stopping Kafka..."
docker stop kafka-gokafka 2>/dev/null && docker rm kafka-gokafka 2>/dev/null

echo "Stopping Zookeeper..."
docker stop zookeeper-gokafka 2>/dev/null && docker rm zookeeper-gokafka 2>/dev/null

echo "Stopping Redis..."
docker stop redis-gokafka 2>/dev/null && docker rm redis-gokafka 2>/dev/null

echo "Stopping PostgreSQL..."
docker stop postgres-gokafka 2>/dev/null && docker rm postgres-gokafka 2>/dev/null

echo "Removing Docker network..."
docker network rm gokafka-network 2>/dev/null

echo "Infrastructure services stopped and cleaned up!"
