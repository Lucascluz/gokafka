# GoKafka Scripts

This directory contains scripts to manage the GoKafka system infrastructure and services.

## Quick Start

To start the entire system:
```bash
./scripts/start-all.sh
```

To stop the entire system:
```bash
./scripts/stop-all.sh
```

## Individual Scripts

### Infrastructure Management

- **start-infrastructure.sh** - Starts all Docker containers (PostgreSQL, Redis, Kafka, Zookeeper)
- **stop-infrastructure.sh** - Stops and removes all Docker containers

### Service Management

- **start-services.sh** - Starts all Go services (API Gateway, User Service)
- **stop-services.sh** - Stops all Go services

## Infrastructure Services

When running, the following services will be available:

- **PostgreSQL**: `localhost:5432`
  - Database: `gokafka`
  - Username: `postgres`
  - Password: `postgres`

- **Redis**: `localhost:6379`

- **Kafka**: `localhost:9092`

- **Zookeeper**: `localhost:2181`

## Application Services

- **API Gateway**: `http://localhost:8080`
  - Test endpoint: `http://localhost:8080/test`

- **User Service**: Kafka consumer (no HTTP endpoint)

## Logs

Service logs are stored in the `logs/` directory:
- `logs/api-gateway.log`
- `logs/user-service.log`

## Requirements

- Docker
- Go 1.22+
- Linux/macOS with bash

## Notes

- All scripts must be run from the project root directory
- Scripts will automatically create necessary directories
- Docker containers use a dedicated network `gokafka-network`
- Services are automatically restarted if already running
