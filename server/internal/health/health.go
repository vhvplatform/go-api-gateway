package health

import (
	"context"
	"time"
)

// HealthChecker manages health checks for various services
type HealthChecker struct {
	checks map[string]HealthCheck
}

// HealthCheck is a function that checks the health of a service
type HealthCheck func(ctx context.Context) error

// HealthStatus represents the overall health status
type HealthStatus struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// NewHealthChecker creates a new health checker
func NewHealthChecker() *HealthChecker {
	return &HealthChecker{
		checks: make(map[string]HealthCheck),
	}
}

// RegisterCheck registers a health check for a service
func (h *HealthChecker) RegisterCheck(name string, check HealthCheck) {
	h.checks[name] = check
}

// CheckAll runs all registered health checks
func (h *HealthChecker) CheckAll(ctx context.Context) HealthStatus {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	status := HealthStatus{
		Status:   "healthy",
		Services: make(map[string]string),
	}

	for name, check := range h.checks {
		if err := check(ctx); err != nil {
			status.Services[name] = "unhealthy: " + err.Error()
			status.Status = "degraded"
		} else {
			status.Services[name] = "healthy"
		}
	}

	return status
}
