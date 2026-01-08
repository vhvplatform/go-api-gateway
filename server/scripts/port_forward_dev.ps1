<#
.SYNOPSIS
    Helper to port-forward Dev services to localhost.
.DESCRIPTION
    Runs kubectl port-forward for shared services.
    Requires kubectl to be configured for the Dev cluster.
#>

Write-Host "Starting Port Forwards for Dev Environment..." -ForegroundColor Cyan

# Check kubectl
if (-not (Get-Command "kubectl" -ErrorAction SilentlyContinue)) {
    Write-Error "kubectl not found! Please install it and configure your context."
    exit 1
}

# Run in background jobs
Write-Host "Forwarding Auth Service (50051 -> 50051)..."
Start-Job -ScriptBlock { kubectl port-forward svc/auth-service 50051:50051 }

Write-Host "Forwarding Tenant Service (50053 -> 50051)..."
Start-Job -ScriptBlock { kubectl port-forward svc/tenant-service 50053:50051 } # Mapping different local port if needed, or same.

Write-Host "Forwarding Notification Service (8084 -> 80)..."
Start-Job -ScriptBlock { kubectl port-forward svc/notification-service 8084:80 }

Write-Host "Forwarding Postgres (5432 -> 5432)..."
Start-Job -ScriptBlock { kubectl port-forward svc/postgres 5432:5432 }

Write-Host "Forwarding Redis (6379 -> 6379)..."
Start-Job -ScriptBlock { kubectl port-forward svc/redis 6379:6379 }

Write-Host "Port forwarding started in background. Close this window to stop (or Stop-Job)."
Read-Host "Press Enter to exit..."
