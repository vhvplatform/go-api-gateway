# VHV Platform - API Gateway

This is a microservices-based platform with multiple components.

## Repository Structure

This repository is organized into the following main directories:

### 📁 `/server` - Backend (Golang)
The Golang-based API Gateway microservice providing backend services.
- **Technology**: Go 1.25.5
- **Purpose**: API Gateway with authentication, routing, rate limiting, circuit breaker, and more
- **Documentation**: See [server/README.md](docs/README.md)

### 📁 `/client` - Frontend (ReactJS)
The ReactJS-based frontend microservice for web applications.
- **Technology**: ReactJS
- **Purpose**: Web user interface
- **Status**: Coming soon

### 📁 `/flutter` - Mobile App (Flutter)
The Flutter-based mobile application.
- **Technology**: Flutter
- **Purpose**: Cross-platform mobile application (iOS/Android)
- **Status**: Coming soon

### 📁 `/docs` - Documentation
Comprehensive project documentation including:
- Architecture diagrams
- API specifications
- Setup guides (Windows, Linux, macOS)
- Contributing guidelines
- Troubleshooting guides
- Examples and tutorials

## Quick Start

### Backend (Server)
```bash
cd server
make build
make run
```

See [docs/README.md](docs/README.md) for detailed backend documentation.

### Frontend (Client)
Coming soon - ReactJS frontend under development.

### Mobile App (Flutter)
Coming soon - Flutter mobile app under development.

## Documentation

Full documentation is available in the [docs/](docs/) directory:
- **[Main Documentation](docs/README.md)** - Complete backend API Gateway documentation
- **[Contributing Guide](docs/CONTRIBUTING.md)** - Development guidelines
- **[Troubleshooting](docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[Windows Setup](docs/WINDOWS_SETUP.md)** - Windows development setup
- **[Architecture Diagrams](docs/diagrams/)** - System architecture (PlantUML)
- **[Examples](docs/examples/)** - Usage examples

## Getting Help

- Check the [Troubleshooting Guide](docs/TROUBLESHOOTING.md)
- Review the [Contributing Guide](docs/CONTRIBUTING.md)
- Check existing issues on GitHub

## License

See main repository LICENSE file.
