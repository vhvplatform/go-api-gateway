# VHV Platform - API Gateway Repository

This repository contains the microservices for the VHV Platform API Gateway project.

## Repository Structure

The repository is organized into the following main directories:

### 📁 `/server` - Backend (Golang)
Production-ready API Gateway service built with Go 1.25.5, featuring:
- Circuit breaker, rate limiting, distributed tracing
- JWT authentication and multi-tenancy support
- Prometheus metrics and comprehensive monitoring
- 96.4% test coverage

See [server/README.md](server/README.md) for detailed documentation.

### 📁 `/client` - Frontend (ReactJS)
Microservice code for the frontend application using ReactJS.

See [client/README.md](client/README.md) for setup and usage.

### 📁 `/flutter` - Mobile App (Flutter)
Mobile application code built with Flutter.

See [flutter/README.md](flutter/README.md) for setup and usage.

### 📁 `/docs` - Documentation
Comprehensive project documentation including:
- **Architecture Diagrams** (`docs/diagrams/`) - PlantUML system architecture
- **Development Guides** (`docs/development/`) - Contributing, setup, and testing guides
- **Troubleshooting** (`docs/TROUBLESHOOTING.md`) - Common issues and solutions
- **Upgrade Information** (`docs/UPGRADE_SUMMARY.md`) - Version upgrade notes

## Quick Start

### Backend (Server)
```bash
cd server
make build
make run
```

### Frontend (Client)
```bash
cd client
# Instructions to be added
```

### Mobile App (Flutter)
```bash
cd flutter
# Instructions to be added
```

## Documentation

- [Server Documentation](server/README.md) - Complete backend API Gateway documentation
- [Contributing Guide](docs/development/CONTRIBUTING.md) - Development and contribution guidelines
- [Windows Setup](docs/development/WINDOWS_SETUP.md) - Windows development environment setup
- [Troubleshooting](docs/TROUBLESHOOTING.md) - Common issues and solutions
- [Architecture Diagrams](docs/diagrams/) - Visual system architecture

## Development

This project follows the architectural standards defined in [go-infrastructure](https://github.com/vhvplatform/go-infrastructure).

For detailed development guidelines, see [docs/development/CONTRIBUTING.md](docs/development/CONTRIBUTING.md).

## Features

- **Backend API Gateway**: Production-ready with circuit breaker, rate limiting, tracing
- **Multi-tenancy**: Tenant isolation and context management
- **Security**: JWT authentication, request validation, rate limiting
- **Observability**: Prometheus metrics, distributed tracing with Jaeger
- **High Availability**: Circuit breaker, health checks, graceful shutdown

## License

See LICENSE file for details.

---

**Note**: This repository contains the complete microservices stack for the VHV Platform API Gateway project.
