#!/bin/bash

# Script to start all required infrastructure services (Kafka, Redis, PostgreSQL)

echo "Starting infrastructure services..."

# Create a Docker network for the services
echo "Creating Docker network 'gokafka-network'..."
docker network create gokafka-network 2>/dev/null || echo "Network already exists"

# Start PostgreSQL
echo "Starting PostgreSQL..."
docker run -d \
  --name postgres-gokafka \
  --network gokafka-network \
  -e POSTGRES_DB=gokafka \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 5432:5432 \
  postgres:15-alpine

# Start Redis
echo "Starting Redis..."
docker run -d \
  --name redis-gokafka \
  --network gokafka-network \
  -p 6379:6379 \
  redis:7-alpine

# Start Zookeeper (required for Kafka)
echo "Starting Zookeeper..."
docker run -d \
  --name zookeeper-gokafka \
  --network gokafka-network \
  -e ZOOKEEPER_CLIENT_PORT=2181 \
  -e ZOOKEEPER_TICK_TIME=2000 \
  -p 2181:2181 \
  confluentinc/cp-zookeeper:latest

# Wait a bit for Zookeeper to start
echo "Waiting for Zookeeper to start..."
sleep 10

# Start Kafka
echo "Starting Kafka..."
docker run -d \
  --name kafka-gokafka \
  --network gokafka-network \
  -e KAFKA_BROKER_ID=1 \
  -e KAFKA_ZOOKEEPER_CONNECT=zookeeper-gokafka:2181 \
  -e KAFKA_ADVERTISED_LISTENERS=PLAINTEXT://localhost:9092 \
  -e KAFKA_LISTENER_SECURITY_PROTOCOL_MAP=PLAINTEXT:PLAINTEXT \
  -e KAFKA_INTER_BROKER_LISTENER_NAME=PLAINTEXT \
  -e KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR=1 \
  -e KAFKA_AUTO_CREATE_TOPICS_ENABLE=true \
  -p 9092:9092 \
  confluentinc/cp-kafka:7.0.1

echo "Waiting for services to be ready..."
sleep 15

echo "Checking service status..."
echo "PostgreSQL: $(docker ps --format 'table {{.Names}}\t{{.Status}}' | grep postgres-gokafka)"
echo "Redis: $(docker ps --format 'table {{.Names}}\t{{.Status}}' | grep redis-gokafka)"
echo "Zookeeper: $(docker ps --format 'table {{.Names}}\t{{.Status}}' | grep zookeeper-gokafka)"
echo "Kafka: $(docker ps --format 'table {{.Names}}\t{{.Status}}' | grep kafka-gokafka)"

echo ""
echo "Infrastructure services started successfully!"
echo "PostgreSQL: localhost:5432 (user: postgres, password: postgres, db: gokafka)"
echo "Redis: localhost:6379"
echo "Kafka: localhost:9092"
echo "Zookeeper: localhost:2181"
