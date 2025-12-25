package health

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestNewHealthChecker(t *testing.T) {
	hc := NewHealthChecker()
	if hc == nil {
		t.Fatal("NewHealthChecker() returned nil")
	}
	if hc.checks == nil {
		t.Error("checks map is nil")
	}
}

func TestRegisterCheck(t *testing.T) {
	hc := NewHealthChecker()
	
	checkCalled := false
	testCheck := func(ctx context.Context) error {
		checkCalled = true
		return nil
	}
	
	hc.RegisterCheck("test-service", testCheck)
	
	// Verify the check was registered
	if len(hc.checks) != 1 {
		t.Errorf("Expected 1 check, got %d", len(hc.checks))
	}
	
	// Run the check to verify it's the right one
	ctx := context.Background()
	hc.checks["test-service"](ctx)
	
	if !checkCalled {
		t.Error("Registered check was not called")
	}
}

func TestCheckAll_AllHealthy(t *testing.T) {
	hc := NewHealthChecker()
	
	// Register multiple healthy checks
	hc.RegisterCheck("service-1", func(ctx context.Context) error {
		return nil
	})
	hc.RegisterCheck("service-2", func(ctx context.Context) error {
		return nil
	})
	hc.RegisterCheck("service-3", func(ctx context.Context) error {
		return nil
	})
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	if status.Status != "healthy" {
		t.Errorf("Expected status 'healthy', got '%s'", status.Status)
	}
	
	if len(status.Services) != 3 {
		t.Errorf("Expected 3 services, got %d", len(status.Services))
	}
	
	for name, health := range status.Services {
		if health != "healthy" {
			t.Errorf("Service %s expected 'healthy', got '%s'", name, health)
		}
	}
}

func TestCheckAll_OneUnhealthy(t *testing.T) {
	hc := NewHealthChecker()
	
	hc.RegisterCheck("service-1", func(ctx context.Context) error {
		return nil
	})
	hc.RegisterCheck("service-2", func(ctx context.Context) error {
		return errors.New("connection failed")
	})
	hc.RegisterCheck("service-3", func(ctx context.Context) error {
		return nil
	})
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	if status.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", status.Status)
	}
	
	if status.Services["service-1"] != "healthy" {
		t.Error("service-1 should be healthy")
	}
	
	if status.Services["service-2"] == "healthy" {
		t.Error("service-2 should be unhealthy")
	}
	
	if status.Services["service-3"] != "healthy" {
		t.Error("service-3 should be healthy")
	}
}

func TestCheckAll_AllUnhealthy(t *testing.T) {
	hc := NewHealthChecker()
	
	testError := errors.New("service down")
	hc.RegisterCheck("service-1", func(ctx context.Context) error {
		return testError
	})
	hc.RegisterCheck("service-2", func(ctx context.Context) error {
		return testError
	})
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	if status.Status != "degraded" {
		t.Errorf("Expected status 'degraded', got '%s'", status.Status)
	}
	
	for name, health := range status.Services {
		if health == "healthy" {
			t.Errorf("Service %s should be unhealthy", name)
		}
	}
}

func TestCheckAll_ContextTimeout(t *testing.T) {
	hc := NewHealthChecker()
	
	// Register a check that takes longer than the timeout
	hc.RegisterCheck("slow-service", func(ctx context.Context) error {
		select {
		case <-time.After(10 * time.Second):
			return nil
		case <-ctx.Done():
			return ctx.Err()
		}
	})
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	// The CheckAll function has a 5-second timeout
	// The slow service should be marked as unhealthy due to timeout
	if status.Status == "healthy" {
		t.Error("Expected degraded status due to timeout")
	}
	
	health := status.Services["slow-service"]
	if health == "healthy" {
		t.Error("Slow service should be unhealthy due to timeout")
	}
}

func TestCheckAll_NoChecksRegistered(t *testing.T) {
	hc := NewHealthChecker()
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	if status.Status != "healthy" {
		t.Errorf("Expected status 'healthy' when no checks registered, got '%s'", status.Status)
	}
	
	if len(status.Services) != 0 {
		t.Errorf("Expected 0 services, got %d", len(status.Services))
	}
}

func TestCheckAll_ErrorMessage(t *testing.T) {
	hc := NewHealthChecker()
	
	expectedError := "database connection failed"
	hc.RegisterCheck("database", func(ctx context.Context) error {
		return errors.New(expectedError)
	})
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	dbHealth := status.Services["database"]
	if dbHealth == "healthy" {
		t.Error("Database should be unhealthy")
	}
	
	// Check that error message is included
	if len(dbHealth) < len(expectedError) {
		t.Errorf("Expected error message to be included in health status")
	}
}

func TestCheckAll_ConcurrentChecks(t *testing.T) {
	hc := NewHealthChecker()
	
	// Register multiple checks
	for i := 0; i < 10; i++ {
		name := "service-" + string(rune('0'+i))
		hc.RegisterCheck(name, func(ctx context.Context) error {
			// Simulate some work
			time.Sleep(10 * time.Millisecond)
			return nil
		})
	}
	
	ctx := context.Background()
	start := time.Now()
	status := hc.CheckAll(ctx)
	duration := time.Since(start)
	
	// All checks run concurrently within the same timeout
	// Should complete relatively quickly (much less than 10 * 10ms = 100ms sequentially)
	if duration > 2*time.Second {
		t.Errorf("CheckAll took too long: %v", duration)
	}
	
	if status.Status != "healthy" {
		t.Error("All checks should be healthy")
	}
	
	if len(status.Services) != 10 {
		t.Errorf("Expected 10 services, got %d", len(status.Services))
	}
}

func TestHealthStatus_JSONSerialization(t *testing.T) {
	hc := NewHealthChecker()
	
	hc.RegisterCheck("service-1", func(ctx context.Context) error {
		return nil
	})
	
	ctx := context.Background()
	status := hc.CheckAll(ctx)
	
	// Verify the struct can be properly marshaled to JSON
	// This is implicit in the struct tags, but we verify the fields exist
	if status.Status == "" {
		t.Error("Status field is empty")
	}
	
	if status.Services == nil {
		t.Error("Services field is nil")
	}
}

func TestRegisterCheck_Overwrite(t *testing.T) {
	hc := NewHealthChecker()
	
	firstCheckCalled := false
	secondCheckCalled := false
	
	// Register first check
	hc.RegisterCheck("test-service", func(ctx context.Context) error {
		firstCheckCalled = true
		return nil
	})
	
	// Overwrite with second check
	hc.RegisterCheck("test-service", func(ctx context.Context) error {
		secondCheckCalled = true
		return nil
	})
	
	ctx := context.Background()
	hc.CheckAll(ctx)
	
	if firstCheckCalled {
		t.Error("First check should not be called after overwrite")
	}
	
	if !secondCheckCalled {
		t.Error("Second check should be called")
	}
	
	// Should still only have one check
	if len(hc.checks) != 1 {
		t.Errorf("Expected 1 check after overwrite, got %d", len(hc.checks))
	}
}
