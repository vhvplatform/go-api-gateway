package circuitbreaker

import (
	"context"
	"sync"
	"time"

	"github.com/sony/gobreaker"
)

// CircuitBreaker manages circuit breakers for different services
type CircuitBreaker struct {
	breakers map[string]*gobreaker.CircuitBreaker
	mu       sync.RWMutex
}

// NewCircuitBreaker creates a new circuit breaker manager
func NewCircuitBreaker() *CircuitBreaker {
	return &CircuitBreaker{
		breakers: make(map[string]*gobreaker.CircuitBreaker),
	}
}

// GetBreaker returns a circuit breaker for the given service name
func (cb *CircuitBreaker) GetBreaker(name string) *gobreaker.CircuitBreaker {
	cb.mu.RLock()
	breaker, exists := cb.breakers[name]
	cb.mu.RUnlock()

	if exists {
		return breaker
	}

	cb.mu.Lock()
	defer cb.mu.Unlock()

	// Double-check after acquiring write lock
	if breaker, exists := cb.breakers[name]; exists {
		return breaker
	}

	settings := gobreaker.Settings{
		Name:        name,
		MaxRequests: 3,
		Interval:    time.Minute,
		Timeout:     30 * time.Second,
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6
		},
	}

	breaker = gobreaker.NewCircuitBreaker(settings)
	cb.breakers[name] = breaker
	return breaker
}

// Execute wraps a function call with circuit breaker protection
func (cb *CircuitBreaker) Execute(name string, fn func() (interface{}, error)) (interface{}, error) {
	breaker := cb.GetBreaker(name)
	return breaker.Execute(fn)
}

// ExecuteContext wraps a context-aware function call with circuit breaker protection
func (cb *CircuitBreaker) ExecuteContext(ctx context.Context, name string, fn func() (interface{}, error)) (interface{}, error) {
	// Check if context is already cancelled
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	breaker := cb.GetBreaker(name)
	return breaker.Execute(fn)
}
