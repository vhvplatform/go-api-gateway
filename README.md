# API Gateway Service

Production-ready API Gateway with advanced features including circuit breaker, rate limiting, distributed tracing, and comprehensive monitoring.

## Features

### Core Functionality
- **API Routing**: Unified entry point for all microservices
- **Authentication & Authorization**: JWT token validation
- **Multi-tenancy Support**: Tenant isolation and context management
- **API Versioning**: Support for API v1 (v2 ready)

### Production Features

#### 1. Rate Limiting
- **Per-IP Rate Limiting**: Configurable requests per second and burst limits
- **Per-Tenant Rate Limiting**: Separate limits for different tenants
- **Memory-Efficient**: Automatic cleanup of inactive rate limiters (10-minute inactivity)
- **Configurable**: Set via environment variables

#### 2. Circuit Breaker
- **Automatic Failure Detection**: Monitors service health
- **Graceful Degradation**: Prevents cascade failures
- **Configurable Thresholds**: Customizable failure ratios and timeouts
- **Per-Service Breakers**: Independent circuit breakers for each downstream service

#### 3. Distributed Tracing
- **OpenTelemetry Integration**: Standards-based tracing
- **Jaeger Support**: Visual trace analysis
- **Request Correlation**: Track requests across services
- **Performance Insights**: Identify bottlenecks

#### 4. Metrics & Monitoring
- **Prometheus Metrics**: Standards-compliant metrics endpoint
- **Request Metrics**: Total requests, duration, active requests
- **Circuit Breaker Metrics**: Monitor breaker states
- **Custom Metrics**: Easy to add application-specific metrics

#### 5. Redis Caching
- **Response Caching**: Reduce backend load
- **Configurable TTL**: Set expiration per cache entry
- **Optional**: Can run without Redis

#### 6. Enhanced Error Handling
- **Structured Errors**: Consistent error response format
- **Correlation IDs**: Track errors across services
- **Panic Recovery**: Automatic recovery with logging
- **Detailed Error Context**: Timestamps and trace IDs

#### 7. Request Validation
- **Content-Type Validation**: Ensures proper headers
- **Request Size Limits**: Prevents oversized requests (configurable, default 10MB)
- **Timeout Protection**: Prevents long-running requests (30s default)

#### 8. Compression
- **Gzip Compression**: Reduces bandwidth usage
- **Automatic**: Enabled for all responses

#### 9. Health Checks
- **Liveness Check**: `/health` endpoint with service status
- **Readiness Check**: `/ready` endpoint
- **Dependency Checks**: Monitor downstream service health

## Configuration

### Environment Variables

```bash
# Server Configuration
API_GATEWAY_PORT=8080                    # HTTP server port

# Service URLs
AUTH_SERVICE_URL=auth-service:50051      # Auth service gRPC endpoint
USER_SERVICE_URL=user-service:50052      # User service gRPC endpoint
TENANT_SERVICE_URL=tenant-service:50053  # Tenant service gRPC endpoint
NOTIFICATION_SERVICE_URL=http://notification-service:8084  # Notification HTTP endpoint

# JWT Configuration
JWT_SECRET=your-secret-key               # JWT signing secret

# Rate Limiting
RATE_LIMIT_RPS=100                       # Requests per second (default: 100)
RATE_LIMIT_BURST=200                     # Burst capacity (default: 200)

# Request Limits
MAX_REQUEST_SIZE=10485760                # Max request size in bytes (default: 10MB)

# Optional: Redis Cache
REDIS_URL=redis://redis:6379/0           # Redis connection URL

# Optional: Distributed Tracing
ENABLE_TRACING=true                      # Enable OpenTelemetry tracing
JAEGER_URL=http://jaeger:14268/api/traces  # Jaeger collector endpoint

# Optional: Metrics
ENABLE_METRICS=true                      # Enable Prometheus metrics (default: true)

# Optional: Circuit Breaker
CIRCUIT_BREAKER_ENABLED=true             # Enable circuit breaker (default: true)
```

## Endpoints

### Health & Monitoring
- `GET /health` - Service health status with dependency checks
- `GET /ready` - Readiness probe
- `GET /metrics` - Prometheus metrics endpoint

### API Routes
All application routes are prefixed with `/api/v1`:

#### Authentication (Public)
- `POST /api/v1/auth/register` - User registration
- `POST /api/v1/auth/login` - User login
- `POST /api/v1/auth/refresh` - Refresh access token
- `POST /api/v1/auth/logout` - User logout (requires auth)

#### Users (Protected)
- `POST /api/v1/users` - Create user
- `GET /api/v1/users` - List users
- `GET /api/v1/users/search` - Search users
- `GET /api/v1/users/:id` - Get user by ID
- `PUT /api/v1/users/:id` - Update user
- `DELETE /api/v1/users/:id` - Delete user

#### Tenants (Protected)
- `POST /api/v1/tenants` - Create tenant
- `GET /api/v1/tenants` - List tenants
- `GET /api/v1/tenants/:id` - Get tenant by ID
- `PUT /api/v1/tenants/:id` - Update tenant
- `DELETE /api/v1/tenants/:id` - Delete tenant
- `POST /api/v1/tenants/:id/users` - Add user to tenant
- `DELETE /api/v1/tenants/:id/users/:user_id` - Remove user from tenant

#### Notifications (Protected)
- `POST /api/v1/notifications/email` - Send email notification
- `POST /api/v1/notifications/webhook` - Send webhook notification
- `GET /api/v1/notifications` - List notifications
- `GET /api/v1/notifications/:id` - Get notification by ID

## Building & Running

### Build
```bash
cd services/api-gateway
go build -o api-gateway ./cmd/main.go
```

### Run
```bash
./api-gateway
```

### Docker
```bash
docker build -f services/api-gateway/Dockerfile -t api-gateway .
docker run -p 8080:8080 --env-file .env api-gateway
```

## Architecture

### Middleware Stack (in order)
1. **Recovery**: Panic recovery with logging
2. **Correlation ID**: Adds unique request ID
3. **Logger**: Request/response logging
4. **Metrics**: Prometheus metrics collection
5. **Compression**: Gzip response compression
6. **Validation**: Request validation
7. **Size Limit**: Request size limiting
8. **Timeout**: Request timeout enforcement
9. **CORS**: Cross-origin resource sharing
10. **Rate Limit**: Rate limiting per IP/tenant

### New Internal Packages

```
internal/
├── cache/              # Redis caching implementation
│   └── cache.go
├── circuitbreaker/     # Circuit breaker management
│   └── breaker.go
├── client/             # gRPC clients with retry logic
│   ├── auth_client.go
│   ├── user_client.go
│   └── tenant_client.go
├── errors/             # Structured error responses
│   └── errors.go
├── health/             # Health check management
│   └── health.go
├── metrics/            # Prometheus metrics definitions
│   └── metrics.go
├── middleware/         # HTTP middleware
│   ├── auth.go         # JWT authentication
│   ├── correlation.go  # Request correlation
│   ├── logger.go       # Request logging
│   ├── metrics.go      # Metrics collection
│   ├── rate_limit.go   # Rate limiting (with fix)
│   ├── recovery.go     # Panic recovery
│   ├── timeout.go      # Request timeout
│   └── validation.go   # Request validation
└── tracing/            # Distributed tracing
    └── tracing.go
```

## Monitoring & Observability

### Prometheus Metrics
Access metrics at `http://localhost:8080/metrics`:

- `api_gateway_requests_total` - Total requests by method, endpoint, status
- `api_gateway_request_duration_seconds` - Request duration histogram
- `api_gateway_active_requests` - Currently active requests
- `api_gateway_circuit_breaker_state` - Circuit breaker states

### Distributed Tracing
View traces in Jaeger UI when tracing is enabled:
1. Access Jaeger UI (typically `http://localhost:16686`)
2. Select `api-gateway` service
3. View request traces and performance

### Logs
Structured JSON logs include:
- Request details (method, path, query)
- Response status and latency
- Correlation IDs for request tracking
- Error details with stack traces

## Performance Tuning

### Rate Limiting
- Adjust `RATE_LIMIT_RPS` for higher/lower throughput
- Increase `RATE_LIMIT_BURST` for spiky traffic patterns
- Monitor memory usage with many unique IPs

### Timeouts
- Default 30s request timeout (configurable in code)
- 10s gRPC connection timeout
- 30s graceful shutdown timeout

### Connection Pooling
- gRPC connections are persistent with retry logic
- Automatic reconnection on failure
- 3 retry attempts with exponential backoff

## Testing

```bash
# Run unit tests
go test ./internal/middleware/...

# Run with coverage
go test -cover ./...

# Run integration tests
go test ./tests/integration/...
```

## Security

### Features
- **JWT Validation**: All protected endpoints require valid JWT
- **Request Size Limits**: Prevent DoS attacks
- **Rate Limiting**: Prevent abuse
- **Panic Recovery**: Graceful error handling
- **Input Validation**: Content-Type and request validation

### Best Practices
- Keep JWT_SECRET secure and rotated regularly
- Use HTTPS in production
- Configure appropriate rate limits
- Monitor failed authentication attempts
- Review circuit breaker metrics for service health

## Troubleshooting

### Common Issues

#### High Memory Usage
- Check rate limiter cleanup (should remove inactive limiters after 10 minutes)
- Monitor active connections to downstream services
- Review Redis connection pool if caching enabled

#### Circuit Breaker Opens Frequently
- Check downstream service health
- Review failure thresholds (60% failure rate default)
- Increase timeout values if services are slow

#### Request Timeouts
- Increase timeout middleware duration if needed
- Check downstream service performance
- Review distributed traces for bottlenecks

#### Rate Limiting Too Aggressive
- Increase `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST`
- Consider per-tenant vs per-IP limiting strategy
- Monitor `api_gateway_requests_total` with status 429

## Development

### Adding New Routes
1. Create handler in `internal/handler/`
2. Add route in `internal/router/router.go`
3. Apply appropriate middleware (auth, rate limit)
4. Update this README

### Adding New Middleware
1. Create middleware in `internal/middleware/`
2. Add to middleware stack in `cmd/main.go`
3. Write tests
4. Update documentation

## License

See main repository LICENSE file.
