# API Gateway Examples

This directory contains examples of common use cases for the API Gateway.

## Available Examples

1. [Authentication Flow](authentication-example.md) - Complete authentication workflow
2. [User Management](user-management-example.md) - CRUD operations for users
3. [Multi-Tenancy](multi-tenancy-example.md) - Tenant isolation and management
4. [Rate Limiting](rate-limiting-example.md) - Handling rate limits
5. [Circuit Breaker](circuit-breaker-example.md) - Fault tolerance patterns
6. [Docker Compose Setup](docker-compose-example.yml) - Local development setup

## Quick Start

### Using cURL

```bash
# Health check
curl http://localhost:8080/health

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'

# Use the access token from login response
export TOKEN="eyJhbGciOiJIUzI1NiIs..."

# Get users (authenticated)
curl http://localhost:8080/api/v1/users \
  -H "Authorization: Bearer $TOKEN"
```

### Using HTTPie

```bash
# Login
http POST http://localhost:8080/api/v1/auth/login \
  email=user@example.com password=password123

# Get users
http GET http://localhost:8080/api/v1/users \
  Authorization:"Bearer $TOKEN"
```

### Using Postman

Import the OpenAPI spec from `docs/api/openapi.yaml` into Postman to get a complete collection with all endpoints.

## Environment Variables for Examples

```bash
export API_GATEWAY_URL=http://localhost:8080
export API_GATEWAY_PORT=8080
export AUTH_SERVICE_URL=auth-service:50051
export USER_SERVICE_URL=user-service:50052
export TENANT_SERVICE_URL=tenant-service:50053
export REDIS_URL=redis://localhost:6379/0
export JWT_SECRET=your-secret-key-change-in-production
```

## Testing with Docker Compose

See [docker-compose-example.yml](docker-compose-example.yml) for a complete local setup including all backend services.

```bash
docker-compose -f examples/docker-compose-example.yml up
```
