# Development Setup with Tilt

This document explains how to set up and use the Tilt development environment for the GoKafka microservices project.

## Prerequisites

- Docker
- Kubernetes cluster (kind, minikube, or Docker Desktop)
- Tilt installed
- kubectl configured

## Project Structure

```
gokafka/
├── Tiltfile                 # Tilt configuration
├── k8s/                     # Kubernetes manifests
│   ├── infrastructure/      # Infrastructure services
│   │   ├── postgres.yaml    # PostgreSQL with PVC
│   │   ├── redis.yaml       # Redis cache
│   │   └── kafka.yaml       # Kafka (KRaft mode, no Zookeeper)
│   └── services/           # Application services
│       ├── api-gateway.yaml # API Gateway service
│       └── user-service.yaml # User service
├── services/               # Go microservices
│   ├── api-gateway/        # API Gateway
│   └── user-service/       # User management service
└── shared/                 # Shared libraries
```

## Services Overview

### Infrastructure Services
- **PostgreSQL**: Database with persistent storage
- **Redis**: Cache and session storage
- **Kafka**: Message broker (KRaft mode, no Zookeeper needed)

### Application Services
- **API Gateway** (Port 8080): HTTP API, authentication, routing
- **User Service** (Port 8081): User management, Kafka consumer

## Quick Start

1. **Start Tilt development environment:**
   ```bash
   tilt up
   ```

2. **Access the Tilt UI:**
   Open http://localhost:10350 in your browser

3. **Test the services:**
   - API Gateway: http://localhost:8080/health
   - User Service: http://localhost:8081/health

## API Endpoints

### API Gateway (http://localhost:8080)
- `GET /health` - Health check
- `GET /ready` - Readiness check
- `GET /test` - Test Kafka communication
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/logout` - User logout (requires auth)
- `GET /api/v1/profile` - Get user profile (requires auth)
- `PUT /api/v1/profile` - Update user profile (requires auth)
- `GET /api/v1/admin/users` - List users (admin only)
- `DELETE /api/v1/admin/users/:id` - Delete user (admin only)

### User Service (http://localhost:8081)
- `GET /health` - Health check
- `GET /ready` - Readiness check

## Environment Variables

Both services support the following environment variables:

### Common Variables
- `PORT` - Service port (8080 for api-gateway, 8081 for user-service)
- `KAFKA_BROKERS` - Kafka broker addresses (default: localhost:9092)
- `JWT_SECRET` - JWT signing secret

### Database Variables (User Service)
- `POSTGRES_HOST` - PostgreSQL host (default: localhost)
- `POSTGRES_PORT` - PostgreSQL port (default: 5432)
- `POSTGRES_DB` - Database name (default: gokafka)
- `POSTGRES_USER` - Database user (default: postgres)
- `POSTGRES_PASSWORD` - Database password (default: postgres)

### Cache Variables
- `REDIS_HOST` - Redis host (default: localhost)
- `REDIS_PORT` - Redis port (default: 6379)

## Development Workflow

### Starting Development
```bash
# Start all services
tilt up

# View logs and status
# Open http://localhost:10350 in browser
```

### Making Changes
- Tilt automatically detects file changes
- Docker images are rebuilt automatically
- Kubernetes deployments are updated automatically
- Services restart with new code

### Testing
```bash
# Test API Gateway health
curl http://localhost:8080/health

# Test User Service health
curl http://localhost:8081/health

# Test Kafka communication
curl http://localhost:8080/test
```

### Stopping Development
```bash
# Stop all services and clean up
tilt down
```

## Service Communication

- **HTTP**: External clients → API Gateway
- **Kafka**: API Gateway ↔ User Service
  - Topics: `api-gateway-topic`, `user-service-topic`
  - Message format: JSON with correlation IDs

## Resource Dependencies

Tilt ensures proper startup order:
1. Infrastructure services (postgres, redis, kafka)
2. Application services (user-service, api-gateway)

Each service waits for its dependencies to be ready before starting.

## Troubleshooting

### Common Issues

1. **Build failures**: Check Docker daemon and network connectivity
2. **Kafka connection issues**: Ensure Kafka is running and accessible
3. **Database connection issues**: Check PostgreSQL pod status and logs
4. **Port conflicts**: Ensure ports 8080, 8081 are available

### Checking Logs
```bash
# View all logs in Tilt UI or
kubectl logs -f deployment/api-gateway
kubectl logs -f deployment/user-service
kubectl logs -f deployment/kafka
kubectl logs -f deployment/postgres
kubectl logs -f deployment/redis
```

### Manual Kubernetes Commands
```bash
# Check pod status
kubectl get pods

# Check services
kubectl get svc

# Check deployments
kubectl get deployments

# Restart a service
kubectl rollout restart deployment/api-gateway
```

## Next Steps

The development environment is now ready for:
- Adding new microservices
- Implementing new API endpoints
- Testing Kafka message flows
- Database schema changes
- Adding monitoring and observability

To add a new service:
1. Create service code in `services/new-service/`
2. Add Dockerfile
3. Create Kubernetes manifest in `k8s/services/`
4. Update Tiltfile to include the new service
