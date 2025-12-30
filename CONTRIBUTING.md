# Contributing to API Gateway Service

Thank you for your interest in contributing to the API Gateway Service! This document provides guidelines and instructions for contributing.

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Contribution Guidelines](#contribution-guidelines)
- [Pull Request Process](#pull-request-process)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Documentation](#documentation)

## Code of Conduct

### Our Pledge

We are committed to providing a welcoming and inclusive environment for all contributors, regardless of background or experience level.

### Expected Behavior

- Be respectful and considerate
- Welcome newcomers and help them get started
- Focus on constructive feedback
- Accept responsibility for mistakes

### Unacceptable Behavior

- Harassment, discrimination, or offensive comments
- Trolling or insulting remarks
- Personal or political attacks
- Publishing others' private information

## Getting Started

### Prerequisites

Before contributing, ensure you have:

- **Git** installed and configured
- **Go** 1.25.5 or later
- **Docker** for containerization
- **Make** for build automation (Linux/macOS) or use Windows build scripts
- Access to a Kubernetes cluster for integration testing (optional)

**Windows Users**: See [WINDOWS_SETUP.md](WINDOWS_SETUP.md) for complete Windows setup instructions including PowerShell and CMD scripts.

### Fork and Clone

1. **Fork the repository** on GitHub

2. **Clone your fork**:
```bash
git clone https://github.com/YOUR_USERNAME/go-api-gateway.git
cd go-api-gateway
```

3. **Add upstream remote**:
```bash
git remote add upstream https://github.com/vhvplatform/go-api-gateway.git
```

4. **Keep your fork in sync**:
```bash
git fetch upstream
git checkout main
git merge upstream/main
```

## Development Workflow

### 1. Create a Branch

Always create a new branch for your changes:

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/your-bug-fix
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions or updates

### 2. Set Up Development Environment

**Linux/macOS:**
```bash
# Install dependencies
make deps

# Verify setup
make validate
```

**Windows (PowerShell):**
```powershell
# Install dependencies
.\build.ps1 deps

# Verify setup
.\build.ps1 validate
```

**Windows (CMD):**
```cmd
# Install dependencies
build.bat deps

# Verify setup
build.bat validate
```

### 3. Make Changes

- Follow the [Coding Standards](#coding-standards) below
- Keep changes focused and atomic
- Write clear, descriptive commit messages
- Test your changes thoroughly

### 4. Test Your Changes

**Linux/macOS:**
```bash
# Run all tests
make test

# Run tests with coverage
make test-coverage

# Check coverage threshold (80%)
make test-coverage-check

# Run linter
make lint

# Run all validation
make validate
```

**Windows (PowerShell):**
```powershell
# Run all tests
.\build.ps1 test

# Run tests with coverage
.\build.ps1 test-coverage

# Check coverage threshold (80%)
.\build.ps1 test-coverage-check

# Run linter
.\build.ps1 lint

# Run all validation
.\build.ps1 validate
```

**Windows (CMD):**
```cmd
# Run all tests
build.bat test

# Run validation
build.bat validate
```

### 5. Commit Your Changes

Write meaningful commit messages:

```bash
git add .
git commit -m "feat: add rate limiting for API endpoints

- Implement token bucket algorithm
- Add per-IP and per-tenant limiting
- Update documentation with configuration options"
```

Commit message format:
- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `style:` - Code style changes (formatting, etc.)
- `refactor:` - Code refactoring
- `test:` - Test changes
- `chore:` - Build/tooling changes

### 6. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a Pull Request on GitHub.

## Contribution Guidelines

### What to Contribute

We welcome contributions in these areas:

#### Core Features
- New middleware components
- Enhanced error handling
- Performance improvements
- Security enhancements

#### Testing
- Unit tests
- Integration tests
- Performance benchmarks
- Test utilities

#### Documentation
- README improvements
- API documentation
- Code comments
- Architecture diagrams
- Usage examples

#### Bug Fixes
- Bug reports with detailed reproduction steps
- Bug fixes with tests

### Before You Start

1. **Check existing issues** - Someone might be working on it
2. **Open an issue** - Discuss major changes before implementing
3. **Review documentation** - Understand the project structure
4. **Test locally** - Verify your changes work

## Pull Request Process

### 1. Before Submitting

- [ ] Code follows project conventions
- [ ] All tests pass (`make test`)
- [ ] Coverage meets threshold (`make test-coverage-check`)
- [ ] Linter passes (`make lint`)
- [ ] Documentation is updated
- [ ] Commit messages are clear
- [ ] Branch is up to date with main

### 2. PR Description

Include in your PR description:

```markdown
## Description
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Documentation update
- [ ] Refactoring
- [ ] Other (describe)

## Changes Made
- List key changes
- Include any breaking changes
- Mention related issues

## Testing
Describe testing performed

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No new warnings generated
- [ ] Tests added/updated
- [ ] Coverage maintained/improved
```

### 3. Review Process

- Reviewers will be automatically assigned
- Address feedback promptly
- Make requested changes
- Re-request review after updates
- Be patient - reviews take time

### 4. Merging

Once approved:
- PRs will be merged by maintainers
- Squash and merge is preferred
- Delete branch after merge

## Coding Standards

### Go Style Guide

Follow standard Go conventions:
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Code Organization

```
internal/
├── cache/           # Redis caching
├── circuitbreaker/  # Circuit breaker management
├── client/          # gRPC clients
├── errors/          # Error handling
├── handler/         # HTTP handlers
├── health/          # Health checks
├── metrics/         # Prometheus metrics
├── middleware/      # HTTP middleware
├── router/          # Route configuration
└── tracing/         # Distributed tracing
```

### Naming Conventions

- **Packages**: Use short, lowercase names
- **Functions**: Use camelCase (exported) or camelCase (unexported)
- **Variables**: Use descriptive names
- **Constants**: Use ALL_CAPS for exported constants
- **Interfaces**: Prefer single-method interfaces with `-er` suffix

### Best Practices

```go
// ✅ Good: Clear function with proper error handling
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    if id == "" {
        return nil, errors.ErrInvalidInput
    }
    
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    
    return user, nil
}

// ✅ Good: Use context for cancellation
func (s *Service) ProcessRequest(ctx context.Context, req *Request) error {
    select {
    case <-ctx.Done():
        return ctx.Err()
    case result := <-s.process(req):
        return result
    }
}

// ✅ Good: Proper error wrapping
return fmt.Errorf("failed to connect to service: %w", err)

// ❌ Bad: Ignoring errors
user, _ := s.GetUser(ctx, id)

// ❌ Bad: Bare returns
func process() (err error) {
    // ...
    return
}
```

### Middleware Development

When creating new middleware:

```go
func MyMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Before request
        startTime := time.Now()
        
        // Process request
        c.Next()
        
        // After request
        duration := time.Since(startTime)
        log.Printf("Request processed in %v", duration)
    }
}
```

### Error Handling

Use the structured error package:

```go
import "github.com/vhvplatform/go-api-gateway/internal/errors"

// Create structured error
err := errors.NewAPIError(
    errors.ErrInvalidInput,
    "Invalid user ID",
    http.StatusBadRequest,
)

// With correlation ID
err.WithCorrelationID(c.GetString("correlation_id"))
```

## Testing

### Unit Tests

```go
func TestMyFunction(t *testing.T) {
    // Arrange
    input := "test"
    expected := "result"
    
    // Act
    result := MyFunction(input)
    
    // Assert
    if result != expected {
        t.Errorf("expected %s, got %s", expected, result)
    }
}
```

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid input", "test@example.com", true, false},
        {"invalid input", "invalid", false, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Coverage Requirements

- Minimum coverage: **80%**
- Focus on business logic
- Test edge cases and error paths
- Use mocks for external dependencies

### Running Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Check coverage threshold
make test-coverage-check

# Run specific package
go test -v ./internal/middleware

# Run with race detector
go test -race ./...
```

## Documentation

### Code Comments

```go
// GetUser retrieves a user by ID from the database.
// It returns ErrNotFound if the user does not exist.
func GetUser(ctx context.Context, id string) (*User, error) {
    // Implementation
}
```

### README Updates

When adding features:
- Update the Features section
- Add configuration examples
- Include usage examples
- Update API documentation

### Architecture Diagrams

Use PlantUML for diagrams in `docs/diagrams/`

## Build and Deploy

### Local Development

```bash
# Build
make build

# Run
make run

# Run with Docker
make docker-build
make docker-run
```

### Make Targets

```bash
make help          # Show all available targets
make build         # Build the application
make test          # Run tests
make lint          # Run linter
make validate      # Run all checks
make docker-build  # Build Docker image
```

## Questions?

If you have questions:

1. Check existing documentation
2. Search closed issues
3. Open a new issue with the `question` label
4. Join our Slack channel: #api-gateway

## Recognition

Contributors will be:
- Listed in release notes
- Mentioned in documentation
- Added to CONTRIBUTORS.md (if significant contributions)

Thank you for contributing! 🎉
