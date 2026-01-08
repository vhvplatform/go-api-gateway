#!/usr/bin/env pwsh
# PowerShell Build Script for Windows
# API Gateway Service Build Script

[CmdletBinding()]
param(
    [Parameter(Position=0)]
    [ValidateSet('help', 'build', 'build-linux', 'test', 'test-coverage', 'test-coverage-check', 
                 'clean', 'deps', 'deps-update', 'deps-verify', 'run', 'fmt', 'vet', 'lint', 
                 'validate', 'version', 'docker-build', 'docker-run')]
    [string]$Command = 'help'
)

# Variables
$AppName = "api-gateway"
$GoVersion = "1.25.5"
$GoCmd = "go"
$GoBase = $PSScriptRoot
$GoBin = Join-Path $GoBase "bin"

# Error handling
$ErrorActionPreference = "Stop"

# Functions
function Write-Info {
    param([string]$Message)
    Write-Host $Message -ForegroundColor Cyan
}

function Write-Success {
    param([string]$Message)
    Write-Host "✓ $Message" -ForegroundColor Green
}

function Write-Error-Custom {
    param([string]$Message)
    Write-Host "✗ $Message" -ForegroundColor Red
}

function Get-Version {
    try {
        $version = git describe --tags --always --dirty 2>$null
        if (-not $version) { $version = "dev" }
        return $version
    }
    catch {
        return "dev"
    }
}

function Get-BuildTime {
    return (Get-Date -Format "yyyy-MM-dd_HH:mm:ss")
}

function Show-Help {
    Write-Host @"
API Gateway Build Script for Windows (PowerShell)

Available commands:
  help                 - Display this help screen
  build                - Build the application for Windows
  build-linux          - Build for Linux (cross-compile)
  test                 - Run tests
  test-coverage        - Run tests with coverage report
  test-coverage-check  - Check test coverage meets minimum threshold (80%)
  clean                - Clean build artifacts
  deps                 - Download dependencies
  deps-update          - Update dependencies
  deps-verify          - Verify dependencies
  run                  - Run the application locally
  fmt                  - Format Go code
  vet                  - Run go vet
  lint                 - Run golangci-lint (if installed)
  validate             - Run all validation checks (fmt, vet, test)
  version              - Display version information
  docker-build         - Build Docker image (requires Docker Desktop)
  docker-run           - Run Docker container locally

Usage: .\build.ps1 [command]
Example: .\build.ps1 build
"@
}

function Build-App {
    Write-Info "Building $AppName for Windows..."
    
    # Create bin directory if it doesn't exist
    if (-not (Test-Path $GoBin)) {
        New-Item -ItemType Directory -Path $GoBin | Out-Null
    }
    
    $version = Get-Version
    $buildTime = Get-BuildTime
    
    Write-Info "Version: $version"
    Write-Info "Build Time: $buildTime"
    
    $ldflags = "-X main.Version=$version -X main.BuildTime=$buildTime"
    
    # Use forward slashes for cross-platform compatibility
    & $GoCmd build -ldflags $ldflags -o "$GoBin/$AppName.exe" ./cmd/main.go
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Build complete: $GoBin\$AppName.exe"
        return $true
    }
    else {
        Write-Error-Custom "Build failed!"
        return $false
    }
}

function Build-Linux {
    Write-Info "Building $AppName for Linux (cross-compile)..."
    
    if (-not (Test-Path $GoBin)) {
        New-Item -ItemType Directory -Path $GoBin | Out-Null
    }
    
    $version = Get-Version
    $buildTime = Get-BuildTime
    
    $ldflags = "-X main.Version=$version -X main.BuildTime=$buildTime"
    
    $env:CGO_ENABLED = "0"
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    
    # Use forward slashes for cross-platform compatibility
    & $GoCmd build -ldflags $ldflags -o "$GoBin/$AppName-linux" ./cmd/main.go
    
    # Restore environment
    Remove-Item Env:\CGO_ENABLED -ErrorAction SilentlyContinue
    Remove-Item Env:\GOOS -ErrorAction SilentlyContinue
    Remove-Item Env:\GOARCH -ErrorAction SilentlyContinue
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Build complete: $GoBin\$AppName-linux"
        return $true
    }
    else {
        Write-Error-Custom "Build failed!"
        return $false
    }
}

function Run-Tests {
    Write-Info "Running tests..."
    
    # Use array to avoid PowerShell glob expansion
    $testArgs = @('test', '-v', '-race', '-coverprofile=coverage.txt', '-covermode=atomic', './...')
    & $GoCmd $testArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Tests complete"
        return $true
    }
    else {
        Write-Error-Custom "Tests failed!"
        return $false
    }
}

function Run-TestCoverage {
    Write-Info "Running tests with coverage report..."
    
    if (-not (Run-Tests)) {
        return $false
    }
    
    Write-Info "Generating coverage report..."
    & $GoCmd tool cover -html=coverage.txt -o coverage.html
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Coverage report: coverage.html"
        
        # Open coverage report in default browser
        if (Test-Path coverage.html) {
            Write-Info "Opening coverage report in browser..."
            Start-Process coverage.html
        }
        return $true
    }
    return $false
}

function Check-TestCoverage {
    Write-Info "Checking coverage threshold..."
    
    if (-not (Run-Tests)) {
        return $false
    }
    
    $coverageOutput = & $GoCmd tool cover -func=coverage.txt | Select-String "total"
    if ($coverageOutput -match "(\d+\.\d+)%") {
        $coverage = [decimal]$matches[1]
        
        if ($coverage -lt 80) {
            Write-Error-Custom "Coverage $coverage% is below 80%"
            return $false
        }
        else {
            Write-Success "Coverage $coverage% meets threshold"
            return $true
        }
    }
    else {
        Write-Error-Custom "Unable to calculate coverage"
        return $false
    }
}

function Clean-Build {
    Write-Info "Cleaning build artifacts..."
    
    if (Test-Path $GoBin) {
        Remove-Item -Recurse -Force $GoBin
    }
    
    Remove-Item -Force coverage.txt -ErrorAction SilentlyContinue
    Remove-Item -Force coverage.html -ErrorAction SilentlyContinue
    Remove-Item -Force test-output.log -ErrorAction SilentlyContinue
    
    Write-Success "Clean complete"
}

function Download-Deps {
    Write-Info "Downloading dependencies..."
    
    & $GoCmd mod download
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Dependencies downloaded"
        return $true
    }
    else {
        Write-Error-Custom "Failed to download dependencies!"
        return $false
    }
}

function Update-Deps {
    Write-Info "Updating dependencies..."
    
    & $GoCmd mod tidy
    $getArgs = @('get', '-u', './...')
    & $GoCmd $getArgs
    & $GoCmd mod tidy
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Dependencies updated"
        return $true
    }
    else {
        Write-Error-Custom "Failed to update dependencies!"
        return $false
    }
}

function Verify-Deps {
    Write-Info "Verifying dependencies..."
    
    & $GoCmd mod verify
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Dependencies verified"
        return $true
    }
    else {
        Write-Error-Custom "Dependency verification failed!"
        return $false
    }
}

function Run-App {
    Write-Info "Running $AppName..."
    & $GoCmd run ./cmd/main.go
}

function Format-Code {
    Write-Info "Formatting code..."
    $fmtArgs = @('fmt', './...')
    & $GoCmd $fmtArgs
    Write-Success "Format complete"
}

function Run-Vet {
    Write-Info "Running go vet..."
    
    $vetArgs = @('vet', './...')
    & $GoCmd $vetArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Vet complete"
        return $true
    }
    else {
        Write-Error-Custom "Vet found issues!"
        return $false
    }
}

function Run-Lint {
    Write-Info "Running golangci-lint..."
    
    # Check if golangci-lint is installed
    $lintCmd = Get-Command golangci-lint -ErrorAction SilentlyContinue
    
    if (-not $lintCmd) {
        Write-Info "golangci-lint not found. Installing..."
        & go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
        
        if ($LASTEXITCODE -ne 0) {
            Write-Error-Custom "Failed to install golangci-lint"
            return $false
        }
    }
    
    $lintArgs = @('run', '--timeout=5m', './...')
    & golangci-lint $lintArgs
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Lint complete"
        return $true
    }
    else {
        Write-Error-Custom "Lint found issues!"
        return $false
    }
}

function Run-Validate {
    Write-Info "Running validation checks..."
    
    Format-Code
    
    if (-not (Run-Vet)) {
        return $false
    }
    
    if (-not (Check-TestCoverage)) {
        return $false
    }
    
    Write-Success "All validation checks passed!"
    return $true
}

function Show-Version {
    Write-Host "Application: $AppName"
    Write-Host "Version: $(Get-Version)"
    Write-Host "Build Time: $(Get-BuildTime)"
    Write-Host "Go Version: $(& $GoCmd version)"
}

function Build-Docker {
    Write-Info "Building Docker image..."
    
    $version = Get-Version
    $dockerImage = "ghcr.io/vhvplatform/$AppName"
    $dockerTag = $version
    
    docker build -t "${dockerImage}:${dockerTag}" .
    
    if ($LASTEXITCODE -eq 0) {
        Write-Success "Docker image built: ${dockerImage}:${dockerTag}"
        return $true
    }
    else {
        Write-Error-Custom "Docker build failed!"
        return $false
    }
}

function Run-Docker {
    Write-Info "Running Docker container..."
    
    $version = Get-Version
    $dockerImage = "ghcr.io/vhvplatform/$AppName"
    $dockerTag = $version
    
    $envFile = ".env.local"
    if (-not (Test-Path $envFile)) {
        $envFile = ".env"
    }
    
    if (Test-Path $envFile) {
        docker run -p 8080:8080 --env-file $envFile "${dockerImage}:${dockerTag}"
    }
    else {
        Write-Info "Warning: No .env or .env.local file found. Running without environment variables."
        docker run -p 8080:8080 "${dockerImage}:${dockerTag}"
    }
}

# Main execution
try {
    switch ($Command) {
        'help' { Show-Help }
        'build' { Build-App }
        'build-linux' { Build-Linux }
        'test' { Run-Tests }
        'test-coverage' { Run-TestCoverage }
        'test-coverage-check' { Check-TestCoverage }
        'clean' { Clean-Build }
        'deps' { Download-Deps }
        'deps-update' { Update-Deps }
        'deps-verify' { Verify-Deps }
        'run' { Run-App }
        'fmt' { Format-Code }
        'vet' { Run-Vet }
        'lint' { Run-Lint }
        'validate' { Run-Validate }
        'version' { Show-Version }
        'docker-build' { Build-Docker }
        'docker-run' { Run-Docker }
        default { Show-Help }
    }
}
catch {
    Write-Error-Custom "An error occurred: $_"
    exit 1
}
