# API Gateway Service

Production-ready API Gateway with advanced features including circuit breaker, rate limiting, distributed tracing, and comprehensive monitoring.

This service follows the architectural standards defined in [go-infrastructure](https://github.com/vhvplatform/go-infrastructure).

[![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg)](https://go.dev/)
[![Test Coverage](https://img.shields.io/badge/coverage-96.4%25-brightgreen.svg)](./coverage.txt)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

## Table of Contents

- [Features](#features)
- [Documentation](#documentation)
- [Configuration](#configuration)
- [Endpoints](#endpoints)
- [Building & Running](#building--running)
- [Architecture](#architecture)
- [Monitoring & Observability](#monitoring--observability)
- [Performance Tuning](#performance-tuning)
- [Testing](#testing)
- [Security](#security)
- [Troubleshooting](#troubleshooting)
- [Contributing](#contributing)

## Documentation

### Quick Links
- **[Contributing Guide](CONTRIBUTING.md)** - Development and contribution guidelines
- **[Windows Setup Guide](WINDOWS_SETUP.md)** - Complete Windows development setup
- **[Troubleshooting Guide](TROUBLESHOOTING.md)** - Common issues and solutions
- **[Architecture Diagrams](docs/diagrams/)** - Visual system architecture (PlantUML)
- **[Examples](examples/)** - Usage examples and Docker Compose setup
- **[API Documentation](docs/api/)** - OpenAPI/Swagger specs
- **[go-infrastructure](https://github.com/vhvplatform/go-infrastructure)** - Infrastructure standards and deployment

### Diagrams
Comprehensive PlantUML diagrams documenting the system:
- [Architecture Overview](docs/diagrams/architecture.puml) - Complete system architecture
- [Request Flow](docs/diagrams/request-flow.puml) - HTTP request lifecycle
- [Authentication](docs/diagrams/authentication.puml) - JWT auth sequences
- [Rate Limiting](docs/diagrams/rate-limiting.puml) - Token bucket algorithm
- [Circuit Breaker](docs/diagrams/circuit-breaker.puml) - Fault tolerance patterns
- [Deployment](docs/diagrams/deployment.puml) - Production topology

### Examples
- [Authentication Flow](examples/authentication-example.md) - Complete auth workflow
- [Docker Compose Setup](examples/docker-compose.yml) - Local development
- [Prometheus Config](examples/prometheus.yml) - Metrics collection

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

This project follows the [go-infrastructure](https://github.com/vhvplatform/go-infrastructure) architectural standards for building, testing, and deployment.

### Prerequisites

- Go 1.25.5 or later
- Docker (for containerization)
- Make (for build automation on Linux/macOS)
- kubectl (for Kubernetes deployment)

**Windows Users**: See [WINDOWS_SETUP.md](WINDOWS_SETUP.md) for detailed Windows development setup including PowerShell and batch scripts.

### Quick Start with Makefile (Linux/macOS)

```bash
# Show all available targets
make help

# Download dependencies
make deps

# Build the application
make build

# Run tests
make test

# Run all validation (fmt, vet, lint, test)
make validate

# Build Docker image
make docker-build
```

### Quick Start on Windows

```powershell
# PowerShell - Show all available commands
.\build.ps1 help

# Download dependencies
.\build.ps1 deps

# Build the application
.\build.ps1 build

# Run tests
.\build.ps1 test

# Run all validation (fmt, vet, test)
.\build.ps1 validate
```

Or using Command Prompt:

```cmd
# CMD - Show all available commands
build.bat help

# Build the application
build.bat build
```

For complete Windows setup instructions, see [WINDOWS_SETUP.md](WINDOWS_SETUP.md).

### Build

**Linux/macOS:**
```bash
# Using Makefile (recommended)
make build

# Or using go directly
go build -o bin/api-gateway ./cmd/main.go
```

**Windows:**
```powershell
# Using PowerShell script (recommended)
.\build.ps1 build

# Or using CMD batch script
build.bat build

# Or using go directly
go build -o bin\api-gateway.exe .\cmd\main.go
```

### Run

**Linux/macOS:**
```bash
# Using Makefile
make run

# Or run the binary directly
./bin/api-gateway
```

**Windows:**
```powershell
# Using PowerShell script
.\build.ps1 run

# Or using CMD batch script
build.bat run

# Or run the binary directly
.\bin\api-gateway.exe
```

### Docker

```bash
# Build Docker image
make docker-build

# Or using docker directly
docker build -t api-gateway:latest .

# Run Docker container
docker run -p 8080:8080 --env-file .env api-gateway:latest
```

### Testing

```bash
# Run all tests
make test

# Run tests with coverage report
make test-coverage

# Check coverage meets threshold (80%)
make test-coverage-check

# Run linter
make lint
```

### Development Workflow

```bash
# Format, vet, lint, test, and build
make validate

# Run local CI pipeline
make ci
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

The project maintains >96% test coverage with comprehensive unit tests. See [CONTRIBUTING.md](CONTRIBUTING.md) for testing guidelines.

```bash
# Run all tests with Makefile (recommended)
make test

# Run tests with coverage report
make test-coverage

# Check coverage threshold (80%)
make test-coverage-check

# View coverage in browser
make test-coverage
open coverage.html

# Run specific package tests
go test -v ./internal/circuitbreaker
go test -v ./internal/health
go test -v ./internal/errors

# Run with race detector
go test -race ./...
```

**Test Coverage**: 96.4% of statements

See the [Testing section in CONTRIBUTING.md](CONTRIBUTING.md#testing) for detailed testing guidelines.

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

For detailed troubleshooting information, see [TROUBLESHOOTING.md](TROUBLESHOOTING.md).

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

This project follows the architectural standards defined in [go-infrastructure](https://github.com/vhvplatform/go-infrastructure). Please refer to the [CONTRIBUTING.md](CONTRIBUTING.md) guide for detailed development guidelines.

### Project Structure

```
.
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── cache/               # Redis caching
│   ├── circuitbreaker/      # Circuit breaker management
│   ├── client/              # gRPC clients
│   ├── errors/              # Error handling
│   ├── handler/             # HTTP handlers
│   ├── health/              # Health checks
│   ├── metrics/             # Prometheus metrics
│   ├── middleware/          # HTTP middleware
│   ├── router/              # Route configuration
│   └── tracing/             # Distributed tracing
├── docs/
│   ├── diagrams/            # PlantUML architecture diagrams
│   └── api/                 # API documentation
├── examples/                # Usage examples
├── Dockerfile               # Container image
├── Makefile                 # Build automation
├── go.mod                   # Go dependencies
├── CONTRIBUTING.md          # Contribution guidelines
└── README.md               # This file
```

### Testing

```bash
# Run all tests (recommended - uses Makefile)
make test

# Run tests with coverage
make test-coverage

# Check coverage threshold
make test-coverage-check

# Or use go directly
go test ./...
go test -coverprofile=coverage.txt -covermode=atomic ./...
go tool cover -html=coverage.txt

# Run specific package tests
go test -v ./internal/circuitbreaker
go test -v ./internal/health
go test -v ./internal/errors
```

**Test Coverage**: 96.4% of statements

### Adding New Routes
1. Create handler in `internal/handler/`
2. Add route in `internal/router/router.go`
3. Apply appropriate middleware (auth, rate limit)
4. Write tests
5. Update this README and API documentation

### Adding New Middleware
1. Create middleware in `internal/middleware/`
2. Add to middleware stack in `cmd/main.go`
3. Write comprehensive tests
4. Update documentation

### Code Quality

The codebase follows Go best practices:
- **Error Handling**: Proper error wrapping and context
- **Context Propagation**: Request context flows through all layers
- **Graceful Shutdown**: 30-second timeout for clean shutdown
- **Concurrency Safety**: All shared resources are properly synchronized
- **Testing**: >96% code coverage with comprehensive unit tests

## Contributing

We welcome contributions! Please read our [Contributing Guide](CONTRIBUTING.md) for details on:

- Development setup and workflow
- Coding standards and best practices
- Testing requirements (minimum 80% coverage)
- Pull request process
- Code review guidelines

### Quick Contribution Steps

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes following our [coding standards](CONTRIBUTING.md#coding-standards)
4. Write tests for your changes (`make test`)
5. Ensure all validation passes (`make validate`)
6. Update documentation as needed
7. Submit a pull request

For detailed guidelines, see [CONTRIBUTING.md](CONTRIBUTING.md).

## Changelog

### v1.2.0 (Latest)
- 🏗️ **Architecture Alignment**: Aligned with [go-infrastructure](https://github.com/vhvplatform/go-infrastructure) standards
- 🔨 **Build System**: Added comprehensive Makefile with 20+ targets
- 📝 **Documentation**: Added CONTRIBUTING.md with detailed development guidelines
- 🐳 **Docker**: Enhanced Dockerfile with multi-stage build and CA certificates
- 🎯 **Build Automation**: Standardized build, test, and deployment processes
- 📦 **Project Structure**: Added .dockerignore and .gitignore for cleaner builds

### v1.1.0
- ✨ Upgraded to Go 1.25.5 (latest stable)
- 📝 Added comprehensive documentation (diagrams, examples, troubleshooting)
- ✅ Achieved 96.4% test coverage with comprehensive unit tests
- 🎨 Added 6 detailed PlantUML architecture diagrams
- 📦 Added Docker Compose setup for local development
- 🔒 Passed security audit (CodeQL) with zero vulnerabilities

### v1.0.0
- Initial release with production features
- Circuit breaker, rate limiting, distributed tracing
- JWT authentication and multi-tenancy
- Prometheus metrics and Redis caching

## License

See main repository LICENSE file.
