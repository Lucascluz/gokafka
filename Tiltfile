# Docker builds for both services with correct build context (root directory)

# Build from the project root to include shared dependencies
docker_build('api-gateway', '.', dockerfile='services/api-gateway/Dockerfile')
docker_build('user-service', '.', dockerfile='services/user-service/Dockerfile')
docker_build('product-service', '.', dockerfile='services/product-service/Dockerfile')

# K8s resource loading for infrastructure and services

# Load infrastructure first
k8s_yaml('k8s/infrastructure/postgres.yaml')
k8s_yaml('k8s/infrastructure/redis.yaml')
k8s_yaml('k8s/infrastructure/kafka.yaml')

# Load services after infrastructure
k8s_yaml('k8s/services/api-gateway.yaml')
k8s_yaml('k8s/services/user-service.yaml')
k8s_yaml('k8s/services/product-service.yaml')

# Resource dependencies and port forwarding
k8s_resource('api-gateway', 
  resource_deps=['postgres', 'redis', 'kafka', 'user-service','product-service'],
  port_forwards='8080:8080'
)

k8s_resource('user-service', 
  resource_deps=['postgres', 'redis', 'kafka'],
  port_forwards='8081:8081'
)

k8s_resource('product-service', 
  resource_deps=['postgres', 'redis', 'kafka'],
  port_forwards='8082:8082'
)