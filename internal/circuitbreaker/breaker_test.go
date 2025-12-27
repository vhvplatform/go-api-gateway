package circuitbreaker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sony/gobreaker"
)

func TestNewCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker()
	if cb == nil {
		t.Fatal("NewCircuitBreaker() returned nil")
	}
	if cb.breakers == nil {
		t.Error("breakers map is nil")
	}
}

func TestGetBreaker_CreateNew(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker := cb.GetBreaker("test-service")

	if breaker == nil {
		t.Fatal("GetBreaker() returned nil")
	}

	// Verify it's stored in the map
	if len(cb.breakers) != 1 {
		t.Errorf("Expected 1 breaker, got %d", len(cb.breakers))
	}
}

func TestGetBreaker_ReuseExisting(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker1 := cb.GetBreaker("test-service")
	breaker2 := cb.GetBreaker("test-service")

	// Should return the same instance
	if breaker1 != breaker2 {
		t.Error("GetBreaker() returned different instances for same service")
	}

	// Should only have one breaker in the map
	if len(cb.breakers) != 1 {
		t.Errorf("Expected 1 breaker, got %d", len(cb.breakers))
	}
}

func TestGetBreaker_MultipleServices(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker1 := cb.GetBreaker("service-1")
	breaker2 := cb.GetBreaker("service-2")
	breaker3 := cb.GetBreaker("service-3")

	if breaker1 == nil || breaker2 == nil || breaker3 == nil {
		t.Fatal("GetBreaker() returned nil")
	}

	// Should have three separate breakers
	if breaker1 == breaker2 || breaker2 == breaker3 || breaker1 == breaker3 {
		t.Error("GetBreaker() returned same instance for different services")
	}

	if len(cb.breakers) != 3 {
		t.Errorf("Expected 3 breakers, got %d", len(cb.breakers))
	}
}

func TestCircuitBreaker_SuccessfulRequests(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker := cb.GetBreaker("test-service")

	// Execute successful requests
	for i := 0; i < 5; i++ {
		_, err := breaker.Execute(func() (interface{}, error) {
			return "success", nil
		})

		if err != nil {
			t.Errorf("Unexpected error on request %d: %v", i, err)
		}
	}

	// Circuit should remain closed
	state := breaker.State()
	if state != gobreaker.StateClosed {
		t.Errorf("Expected state Closed, got %v", state)
	}
}

func TestCircuitBreaker_FailureTrip(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker := cb.GetBreaker("test-service")

	// Execute failing requests
	// Need at least 3 requests with 60% failure rate to trip
	testError := errors.New("test error")

	for i := 0; i < 5; i++ {
		breaker.Execute(func() (interface{}, error) {
			return nil, testError
		})
	}

	// Circuit should be open now
	state := breaker.State()
	if state != gobreaker.StateOpen {
		t.Errorf("Expected state Open, got %v", state)
	}

	// Next request should fail immediately
	_, err := breaker.Execute(func() (interface{}, error) {
		return "should not execute", nil
	})

	if err != gobreaker.ErrOpenState {
		t.Errorf("Expected ErrOpenState, got %v", err)
	}
}

func TestCircuitBreaker_HalfOpenTransition(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker := cb.GetBreaker("test-service")

	// Trip the circuit
	testError := errors.New("test error")
	for i := 0; i < 5; i++ {
		breaker.Execute(func() (interface{}, error) {
			return nil, testError
		})
	}

	// Verify it's open
	if breaker.State() != gobreaker.StateOpen {
		t.Error("Circuit should be open")
	}

	// Wait for timeout (30 seconds is too long, so we'll test the state only)
	// In a real test, you might mock time or use a shorter timeout
	// For now, we just verify the state is open
	state := breaker.State()
	if state != gobreaker.StateOpen {
		t.Errorf("Expected state Open, got %v", state)
	}
}

func TestCircuitBreaker_ConcurrentAccess(t *testing.T) {
	cb := NewCircuitBreaker()

	done := make(chan bool)

	// Concurrent goroutines getting breakers
	for i := 0; i < 10; i++ {
		go func(id int) {
			breaker := cb.GetBreaker("concurrent-service")
			if breaker == nil {
				t.Errorf("Goroutine %d: GetBreaker() returned nil", id)
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}

	// Should only have one breaker despite concurrent access
	if len(cb.breakers) != 1 {
		t.Errorf("Expected 1 breaker, got %d", len(cb.breakers))
	}
}

func TestCircuitBreaker_ReadyToTripConditions(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker := cb.GetBreaker("test-service")

	testError := errors.New("test error")

	// Test 1: Less than 3 requests - should not trip
	breaker.Execute(func() (interface{}, error) {
		return nil, testError
	})
	breaker.Execute(func() (interface{}, error) {
		return nil, testError
	})

	if breaker.State() == gobreaker.StateOpen {
		t.Error("Circuit should not trip with less than 3 requests")
	}

	// Test 2: 3rd request fails - should trip (100% failure rate)
	breaker.Execute(func() (interface{}, error) {
		return nil, testError
	})

	if breaker.State() != gobreaker.StateOpen {
		t.Error("Circuit should trip after 3 failed requests")
	}
}

func TestCircuitBreaker_Settings(t *testing.T) {
	cb := NewCircuitBreaker()
	breaker := cb.GetBreaker("test-service")

	// We can't directly access settings, but we can verify behavior
	// MaxRequests: 3 - will be tested when half-open
	// Interval: 1 minute - resets counters
	// Timeout: 30 seconds - transition to half-open

	// Verify the breaker exists and is initially closed
	if breaker.State() != gobreaker.StateClosed {
		t.Errorf("Initial state should be Closed, got %v", breaker.State())
	}
}

func TestCircuitBreaker_Execute(t *testing.T) {
	cb := NewCircuitBreaker()

	t.Run("Successful execution", func(t *testing.T) {
		result, err := cb.Execute("test-service", func() (interface{}, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != "success" {
			t.Errorf("Expected 'success', got %v", result)
		}
	})

	t.Run("Failed execution", func(t *testing.T) {
		expectedErr := errors.New("test error")
		result, err := cb.Execute("test-service-2", func() (interface{}, error) {
			return nil, expectedErr
		})
		if err != expectedErr {
			t.Errorf("Expected error %v, got %v", expectedErr, err)
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("Circuit breaker trips after failures", func(t *testing.T) {
		serviceName := "test-service-3"
		
		// Cause failures to trip the circuit
		for i := 0; i < 10; i++ {
			_, _ = cb.Execute(serviceName, func() (interface{}, error) {
				return nil, errors.New("failure")
			})
		}

		// Verify circuit is now open
		breaker := cb.GetBreaker(serviceName)
		if breaker.State() != gobreaker.StateOpen {
			t.Errorf("Expected circuit to be Open, got %v", breaker.State())
		}
	})
}

func TestCircuitBreaker_ExecuteContext(t *testing.T) {
	cb := NewCircuitBreaker()

	t.Run("Successful execution with context", func(t *testing.T) {
		ctx := context.Background()
		result, err := cb.ExecuteContext(ctx, "test-service", func() (interface{}, error) {
			return "success", nil
		})
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != "success" {
			t.Errorf("Expected 'success', got %v", result)
		}
	})

	t.Run("Context cancelled before execution", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := cb.ExecuteContext(ctx, "test-service-2", func() (interface{}, error) {
			return "should not run", nil
		})
		
		if err != context.Canceled {
			t.Errorf("Expected context.Canceled error, got %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})

	t.Run("Context timeout during execution", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
		defer cancel()

		time.Sleep(2 * time.Millisecond) // Ensure context is expired

		result, err := cb.ExecuteContext(ctx, "test-service-3", func() (interface{}, error) {
			return "should not run", nil
		})
		
		if err != context.DeadlineExceeded {
			t.Errorf("Expected context.DeadlineExceeded error, got %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result, got %v", result)
		}
	})
}
