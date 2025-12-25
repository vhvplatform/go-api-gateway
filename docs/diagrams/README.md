# API Gateway Architecture Diagrams

This directory contains comprehensive PlantUML diagrams documenting the API Gateway architecture and workflows.

## Diagrams Overview

### 1. Architecture Diagram (`architecture.puml`)
**Purpose**: Shows the overall system architecture and component interactions

**Key Components**:
- API Gateway layer with Gin HTTP server
- Middleware stack (Recovery, Correlation ID, Logger, Metrics, etc.)
- Core services (Circuit Breaker, Redis Cache, Health Checker)
- Backend microservices (Auth, User, Tenant, Notification)
- Observability stack (Prometheus, Jaeger)
- Data layer (MongoDB, Redis)

**Use Cases**:
- Understanding system architecture
- Onboarding new team members
- Architecture reviews
- Documentation for stakeholders

### 2. Request Flow Diagram (`request-flow.puml`)
**Purpose**: Illustrates the complete HTTP request lifecycle through the gateway

**Key Flows**:
- Request processing through middleware chain
- Authentication and authorization checks
- Rate limiting enforcement
- Circuit breaker pattern
- Cache lookup and backend communication
- Response compression and metrics recording

**Use Cases**:
- Debugging request flow issues
- Performance optimization
- Understanding middleware order
- Troubleshooting errors

### 3. Authentication Diagram (`authentication.puml`)
**Purpose**: Documents JWT authentication and authorization workflows

**Key Flows**:
- User registration
- Login and token generation
- Token validation on protected routes
- Token refresh mechanism
- Logout and token revocation
- Multi-tenant context handling

**Use Cases**:
- Implementing authentication clients
- Debugging auth issues
- Security audits
- Understanding token lifecycle

### 4. Rate Limiting Diagram (`rate-limiting.puml`)
**Purpose**: Explains the rate limiting and throttling mechanism

**Key Concepts**:
- Token bucket algorithm
- Per-IP rate limiting
- Concurrent request handling
- Cleanup process for inactive limiters
- Independent limits per IP
- Token refill mechanism

**Use Cases**:
- Configuring rate limits
- Troubleshooting 429 errors
- Understanding memory management
- Capacity planning

### 5. Circuit Breaker Diagram (`circuit-breaker.puml`)
**Purpose**: Shows fault tolerance patterns and state transitions

**Key States**:
- CLOSED: Normal operation
- OPEN: Failing fast (service unhealthy)
- HALF-OPEN: Testing recovery

**Key Concepts**:
- Failure tracking
- State transitions
- Recovery mechanism
- Service isolation
- Metrics export

**Use Cases**:
- Troubleshooting service outages
- Configuring circuit breaker thresholds
- Understanding failure handling
- System reliability design

### 6. Deployment Diagram (`deployment.puml`)
**Purpose**: Documents production deployment topology

**Key Components**:
- Load balancer layer (ALB/Nginx)
- API Gateway cluster (horizontal scaling)
- Backend microservice clusters
- Data layer (MongoDB replica set, Redis cluster)
- Observability stack (Prometheus, Grafana, Jaeger)
- Message queue (RabbitMQ cluster)

**Use Cases**:
- Production deployment planning
- Infrastructure provisioning
- Scaling decisions
- Disaster recovery planning

## Viewing Diagrams

### Online PlantUML Editors
1. **PlantUML Web Server**: https://www.plantuml.com/plantuml/
   - Paste diagram code
   - View rendered diagram
   - Export as PNG/SVG

2. **PlantText**: https://www.planttext.com/
   - Simpler interface
   - Real-time preview

### VS Code Extension
Install the PlantUML extension:
```bash
code --install-extension plantuml.plantuml
```

Then open any `.puml` file and press `Alt+D` to preview.

### Generate PNG/SVG Files

#### Using Docker
```bash
# Generate all diagrams as PNG
docker run --rm -v $(pwd):/data plantuml/plantuml:latest \
  -tpng /data/docs/diagrams/*.puml

# Generate as SVG (scalable)
docker run --rm -v $(pwd):/data plantuml/plantuml:latest \
  -tsvg /data/docs/diagrams/*.puml
```

#### Using PlantUML CLI
```bash
# Install PlantUML (requires Java)
brew install plantuml  # macOS
apt-get install plantuml  # Ubuntu/Debian

# Generate diagrams
plantuml docs/diagrams/*.puml
```

#### Using Node.js
```bash
npm install -g node-plantuml
puml generate docs/diagrams/*.puml -o docs/diagrams/png/
```

### Generate PDF Documentation
```bash
# Generate all diagrams as PDF
docker run --rm -v $(pwd):/data plantuml/plantuml:latest \
  -tpdf /data/docs/diagrams/*.puml
```

## Updating Diagrams

1. **Edit the `.puml` file** using any text editor
2. **Follow PlantUML syntax**: https://plantuml.com/
3. **Test locally** using one of the viewers above
4. **Commit changes** to version control
5. **Regenerate images** if needed for documentation

## Diagram Conventions

### Colors
- **Blue (#E3F2FD)**: API Gateway components
- **Orange (#FFE0B2)**: Backend services
- **Red (#FFCDD2)**: Data stores
- **Default**: Observability and infrastructure

### Labels
- **Solid arrows**: Synchronous calls
- **Dotted arrows**: Asynchronous calls or monitoring
- **Bold text**: Important components
- **Notes**: Additional context and explanations

## Integration with Documentation

These diagrams are referenced in:
- [README.md](../../README.md) - Main documentation
- [TROUBLESHOOTING.md](../../TROUBLESHOOTING.md) - Troubleshooting guide
- [examples/](../../examples/) - Usage examples

## Contributing

When adding new diagrams:
1. Use consistent styling and colors
2. Add descriptive notes for complex flows
3. Update this README with diagram description
4. Reference the diagram in relevant documentation
5. Regenerate PNG/SVG versions

## License

These diagrams are part of the API Gateway documentation and follow the same license as the project.
