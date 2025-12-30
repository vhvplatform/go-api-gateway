# Windows Development Setup Guide

This guide provides detailed instructions for setting up and using the API Gateway service on Windows.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Building the Application](#building-the-application)
- [Running the Application](#running-the-application)
- [Testing](#testing)
- [Development Workflow](#development-workflow)
- [Docker on Windows](#docker-on-windows)
- [IDE Setup](#ide-setup)
- [Troubleshooting](#troubleshooting)
- [Windows-Specific Considerations](#windows-specific-considerations)

---

## Prerequisites

### Required Software

1. **Go 1.25.5 or later**
   - Download from: https://go.dev/dl/
   - During installation, ensure "Add to PATH" is selected
   - Verify installation:
     ```cmd
     go version
     ```

2. **Git for Windows**
   - Download from: https://git-scm.com/download/win
   - Recommended options during installation:
     - Use Git from the Windows Command Prompt
     - Checkout as-is, commit as-is (or Unix-style line endings)
   - Verify installation:
     ```cmd
     git --version
     ```

3. **PowerShell 5.1 or later** (usually pre-installed on Windows 10/11)
   - Check version:
     ```powershell
     $PSVersionTable.PSVersion
     ```
   - For older versions, install PowerShell 7: https://github.com/PowerShell/PowerShell

### Optional Software

4. **Docker Desktop for Windows** (for containerization)
   - Download from: https://www.docker.com/products/docker-desktop
   - System requirements: Windows 10/11 Pro, Enterprise, or Education (64-bit)
   - Enable WSL 2 backend for better performance

5. **Make for Windows** (optional, for using Makefile)
   - Option 1: Install via Chocolatey
     ```cmd
     choco install make
     ```
   - Option 2: Install via Scoop
     ```cmd
     scoop install make
     ```
   - Option 3: Use Git Bash (comes with Git for Windows)
   - Option 4: Use Windows-specific build scripts (recommended, see below)

6. **golangci-lint** (for code linting)
   - Install via Go:
     ```cmd
     go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
     ```

---

## Installation

### 1. Clone the Repository

Open Command Prompt or PowerShell and run:

```cmd
git clone https://github.com/vhvplatform/go-api-gateway.git
cd go-api-gateway
```

### 2. Download Dependencies

Using PowerShell:
```powershell
.\build.ps1 deps
```

Using Command Prompt:
```cmd
build.bat deps
```

Or directly with Go:
```cmd
go mod download
```

---

## Building the Application

### Option 1: Using PowerShell Script (Recommended)

The PowerShell script provides a rich set of features with colored output:

```powershell
# Show all available commands
.\build.ps1 help

# Build the application
.\build.ps1 build

# Build for Linux (cross-compilation)
.\build.ps1 build-linux

# With verbose output
.\build.ps1 build -Verbose
```

The built executable will be in `bin\api-gateway.exe`.

### Option 2: Using Batch Script

For users who prefer CMD or don't have PowerShell:

```cmd
# Show all available commands
build.bat help

# Build the application
build.bat build
```

### Option 3: Using Go Directly

```cmd
# Create bin directory
mkdir bin

# Build
go build -o bin\api-gateway.exe .\cmd\main.go
```

### Option 4: Using Make (if installed)

```bash
# In Git Bash or with Make installed
make build
```

---

## Running the Application

### 1. Set Environment Variables

Create a `.env` file in the project root or set environment variables in PowerShell:

```powershell
# PowerShell
$env:API_GATEWAY_PORT = "8080"
$env:AUTH_SERVICE_URL = "localhost:50051"
$env:USER_SERVICE_URL = "localhost:50052"
$env:TENANT_SERVICE_URL = "localhost:50053"
$env:NOTIFICATION_SERVICE_URL = "http://localhost:8084"
$env:JWT_SECRET = "your-secret-key"
$env:RATE_LIMIT_RPS = "100"
$env:RATE_LIMIT_BURST = "200"
```

Or in Command Prompt:

```cmd
set API_GATEWAY_PORT=8080
set AUTH_SERVICE_URL=localhost:50051
set USER_SERVICE_URL=localhost:50052
set TENANT_SERVICE_URL=localhost:50053
set NOTIFICATION_SERVICE_URL=http://localhost:8084
set JWT_SECRET=your-secret-key
set RATE_LIMIT_RPS=100
set RATE_LIMIT_BURST=200
```

### 2. Run the Application

Using build scripts:

```powershell
# PowerShell
.\build.ps1 run

# Command Prompt
build.bat run
```

Or run the built executable directly:

```cmd
.\bin\api-gateway.exe
```

Or using Go:

```cmd
go run .\cmd\main.go
```

### 3. Verify the Application

Open another terminal and test:

```powershell
# Check health endpoint
curl http://localhost:8080/health

# Or use PowerShell's Invoke-WebRequest
Invoke-WebRequest -Uri http://localhost:8080/health | Select-Object -ExpandProperty Content

# Or use browser
# Navigate to: http://localhost:8080/health
```

---

## Testing

### Run Tests

Using build scripts:

```powershell
# PowerShell - run tests
.\build.ps1 test

# PowerShell - run tests with coverage report
.\build.ps1 test-coverage

# PowerShell - check coverage threshold
.\build.ps1 test-coverage-check

# Command Prompt
build.bat test
```

Using Go directly:

```cmd
# Run all tests
go test .\...

# Run with coverage
go test -v -race -coverprofile=coverage.txt -covermode=atomic .\...

# Generate HTML coverage report
go tool cover -html=coverage.txt -o coverage.html

# Open coverage report in browser
start coverage.html
```

### Run Specific Package Tests

```cmd
# Test circuit breaker
go test -v .\internal\circuitbreaker

# Test health checks
go test -v .\internal\health

# Test with race detector
go test -race .\internal\middleware
```

---

## Development Workflow

### Format, Lint, and Validate

Using PowerShell:

```powershell
# Format code
.\build.ps1 fmt

# Run go vet
.\build.ps1 vet

# Run linter
.\build.ps1 lint

# Run all validation checks (fmt, vet, test)
.\build.ps1 validate
```

Using Command Prompt:

```cmd
# Format code
build.bat fmt

# Run go vet
build.bat vet

# Run validation
build.bat validate
```

### Clean Build Artifacts

```powershell
# PowerShell
.\build.ps1 clean

# Command Prompt
build.bat clean
```

This removes:
- `bin/` directory
- `coverage.txt` and `coverage.html`
- Temporary test files

---

## Docker on Windows

### Prerequisites

- Docker Desktop for Windows installed and running
- WSL 2 backend enabled (recommended)

### Build Docker Image

Using PowerShell:

```powershell
.\build.ps1 docker-build
```

Or using Docker directly:

```cmd
docker build -t api-gateway:latest .
```

### Run Docker Container

Using PowerShell (with .env file):

```powershell
.\build.ps1 docker-run
```

Or using Docker directly:

```cmd
# With environment file
docker run -p 8080:8080 --env-file .env api-gateway:latest

# With individual environment variables
docker run -p 8080:8080 ^
  -e API_GATEWAY_PORT=8080 ^
  -e AUTH_SERVICE_URL=host.docker.internal:50051 ^
  -e USER_SERVICE_URL=host.docker.internal:50052 ^
  -e TENANT_SERVICE_URL=host.docker.internal:50053 ^
  -e JWT_SECRET=your-secret-key ^
  api-gateway:latest
```

**Note**: Use `host.docker.internal` instead of `localhost` to access services running on your Windows host from within Docker containers.

### Docker Compose

```yaml
# docker-compose.yml
version: '3.8'
services:
  api-gateway:
    build: .
    ports:
      - "8080:8080"
    environment:
      - API_GATEWAY_PORT=8080
      - AUTH_SERVICE_URL=host.docker.internal:50051
      - USER_SERVICE_URL=host.docker.internal:50052
      - TENANT_SERVICE_URL=host.docker.internal:50053
      - JWT_SECRET=your-secret-key
```

Run with:

```cmd
docker-compose up -d
```

---

## IDE Setup

### Visual Studio Code

1. **Install Extensions**:
   - Go extension by Go Team at Google
   - GitLens
   - Docker (if using containers)
   - REST Client (for API testing)

2. **Configure Settings** (`.vscode/settings.json`):
   ```json
   {
     "go.useLanguageServer": true,
     "go.lintTool": "golangci-lint",
     "go.lintOnSave": "package",
     "go.formatTool": "gofmt",
     "go.testFlags": ["-v", "-race"],
     "go.coverOnSave": true,
     "editor.formatOnSave": true
   }
   ```

3. **Launch Configuration** (`.vscode/launch.json`):
   ```json
   {
     "version": "0.2.0",
     "configurations": [
       {
         "name": "Launch API Gateway",
         "type": "go",
         "request": "launch",
         "mode": "debug",
         "program": "${workspaceFolder}/cmd/main.go",
         "env": {
           "API_GATEWAY_PORT": "8080",
           "AUTH_SERVICE_URL": "localhost:50051",
           "USER_SERVICE_URL": "localhost:50052",
           "TENANT_SERVICE_URL": "localhost:50053",
           "JWT_SECRET": "your-secret-key"
         }
       }
     ]
   }
   ```

### GoLand / IntelliJ IDEA

1. **Run Configuration**:
   - Go to Run > Edit Configurations
   - Add new "Go Build" configuration
   - Set "Run kind" to "File"
   - Set "Files" to `cmd\main.go`
   - Add environment variables in "Environment" section

2. **Enable Go Modules**:
   - Settings > Go > Go Modules
   - Check "Enable Go modules integration"

---

## Troubleshooting

### Common Windows Issues

#### Issue 1: "go: command not found"

**Solution**:
1. Ensure Go is installed
2. Add Go to PATH:
   - System Properties > Environment Variables
   - Add to Path: `C:\Go\bin` and `%USERPROFILE%\go\bin`
3. Restart terminal

#### Issue 2: "Cannot execute build.ps1 - Execution Policy"

**Error**:
```
.\build.ps1 : File cannot be loaded because running scripts is disabled on this system.
```

**Solution**:
```powershell
# Allow script execution for current user
Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser

# Or run with bypass
powershell -ExecutionPolicy Bypass -File .\build.ps1 build
```

#### Issue 3: Line Ending Issues (CRLF vs LF)

**Symptoms**: Git shows many modified files, or scripts fail to execute

**Solution**:
```cmd
# Configure Git to handle line endings automatically
git config --global core.autocrlf true

# Reset line endings in repository
git rm --cached -r .
git reset --hard
```

#### Issue 4: Port Already in Use

**Error**: `listen tcp :8080: bind: Only one usage of each socket address`

**Solution**:
```powershell
# Find process using port 8080
netstat -ano | findstr :8080

# Kill process (replace PID with actual process ID)
taskkill /PID <PID> /F

# Or change the port
$env:API_GATEWAY_PORT = "8081"
```

#### Issue 5: Long Path Issues

**Error**: `The system cannot find the path specified`

**Solution**:
```powershell
# Enable long paths in Windows (requires admin)
New-ItemProperty -Path "HKLM:\SYSTEM\CurrentControlSet\Control\FileSystem" `
  -Name "LongPathsEnabled" -Value 1 -PropertyType DWORD -Force

# Enable in Git
git config --system core.longpaths true
```

#### Issue 6: Docker "host.docker.internal not found"

**Solution**:
- Ensure Docker Desktop is using WSL 2 backend
- Update Docker Desktop to latest version
- Add to `docker run` command: `--add-host=host.docker.internal:host-gateway`

#### Issue 7: Permission Denied on bin\ Directory

**Solution**:
```cmd
# Close any running instances
taskkill /IM api-gateway.exe /F

# Clean and rebuild
.\build.ps1 clean
.\build.ps1 build
```

---

## Windows-Specific Considerations

### File Paths

- Go accepts both forward slashes (`/`) and backslashes (`\`) in import paths
- Use `filepath.Join()` for cross-platform path handling in code
- Build scripts handle Windows paths automatically

### Environment Variables

- Use `$env:VAR_NAME` in PowerShell
- Use `%VAR_NAME%` in Command Prompt
- Create `.env` file for persistent configuration

### Line Endings

- The repository uses Unix-style (LF) line endings
- Git on Windows handles conversion automatically with `core.autocrlf`
- Recommended setting: `git config --global core.autocrlf true`

### Case Sensitivity

- Windows filesystems are case-insensitive by default
- Go is case-sensitive
- Be careful with package names and imports

### Networking

- Docker containers need `host.docker.internal` to access Windows host
- Firewall may block ports - add exceptions if needed:
  ```powershell
  # Add firewall rule for port 8080
  New-NetFirewallRule -DisplayName "API Gateway" -Direction Inbound -LocalPort 8080 -Protocol TCP -Action Allow
  ```

### Performance

- WSL 2 provides better performance for Docker
- Consider using WSL 2 for development if experiencing performance issues
- Antivirus may slow down builds - add project directory to exclusions

### Tools

- PowerShell Core (PowerShell 7+) recommended for better cross-platform compatibility
- Windows Terminal provides better experience than cmd.exe
- Git Bash provides Unix-like environment

---

## Quick Reference

### PowerShell Commands

```powershell
# Build and run
.\build.ps1 build
.\build.ps1 run

# Testing
.\build.ps1 test
.\build.ps1 test-coverage

# Validation
.\build.ps1 validate

# Docker
.\build.ps1 docker-build
.\build.ps1 docker-run

# Cleanup
.\build.ps1 clean
```

### CMD/Batch Commands

```cmd
# Build and run
build.bat build
build.bat run

# Testing
build.bat test

# Validation
build.bat validate

# Cleanup
build.bat clean
```

### Direct Go Commands

```cmd
# Build
go build -o bin\api-gateway.exe .\cmd\main.go

# Run
go run .\cmd\main.go

# Test
go test .\...

# Format
go fmt .\...

# Vet
go vet .\...
```

---

## Additional Resources

- [Go on Windows Documentation](https://golang.org/doc/install/windows)
- [Docker Desktop for Windows](https://docs.docker.com/desktop/windows/)
- [Git for Windows](https://git-scm.com/download/win)
- [PowerShell Documentation](https://docs.microsoft.com/en-us/powershell/)
- [Visual Studio Code Go Extension](https://marketplace.visualstudio.com/items?itemName=golang.go)

---

## Getting Help

If you encounter issues not covered in this guide:

1. Check the main [TROUBLESHOOTING.md](TROUBLESHOOTING.md)
2. Check the [README.md](README.md) for general documentation
3. Search [GitHub Issues](https://github.com/vhvplatform/go-api-gateway/issues)
4. Create a new issue with:
   - Windows version
   - Go version (`go version`)
   - PowerShell version (`$PSVersionTable.PSVersion`)
   - Detailed error messages
   - Steps to reproduce

---

**Note**: For production deployments on Windows Server, consider using Windows containers or deploying to a Linux-based container orchestration platform like Kubernetes.
