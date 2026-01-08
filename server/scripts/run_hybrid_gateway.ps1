<#
.SYNOPSIS
    Runs the API Gateway in Hybrid Development Mode.
.DESCRIPTION
    This script sets up environment variables to point specific services to localhost (Dev via Port-Forward)
    and others to local instances.
    
    By default, it assumes:
    - Auth Service: Dev (localhost:50051 via port-forward)
    - Tenant Service: Dev (localhost:50053 via port-forward)
    - User Service: Local (localhost:50052 running locally) <- The service you are developing
#>

Write-Host "Starting API Gateway in HYBRID Mode..." -ForegroundColor Cyan

# 1. Gateway Config
$env:API_GATEWAY_PORT = "8080"
$env:JWT_SECRET = "dev-shared-secret" # Must match Dev environment secret
$env:MTLS_ENABLED = "false"           # Disable mTLS for local dev simplicity
$env:ENABLE_METRICS = "true"

# 2. Service Configuration
# Assume we are developing User Service locally, so we point to it directly.
# The others are port-forwarded from Dev.

# Auth Service (Dev - Port Forwarded)
$env:AUTH_SERVICE_URL = "localhost:50051"

# Tenant Service (Dev - Port Forwarded)
$env:TENANT_SERVICE_URL = "localhost:50053"

# User Service (LOCAL - Running on your machine)
$env:USER_SERVICE_URL = "localhost:50052"

# Notification Service (Dev - Port Forwarded)
$env:NOTIFICATION_SERVICE_URL = "http://localhost:8084"

# 3. Cache
$env:CACHE_MAX_COST = "104857600"

Write-Host "Configuration:" -ForegroundColor Gray
Write-Host "  Auth Service   -> $env:AUTH_SERVICE_URL (Expects Port Forward)"
Write-Host "  Tenant Service -> $env:TENANT_SERVICE_URL (Expects Port Forward)"
Write-Host "  User Service   -> $env:USER_SERVICE_URL (Local Instance)"
Write-Host "  Gateway Port   -> $env:API_GATEWAY_PORT"

# 4. Run
go run cmd/main.go
