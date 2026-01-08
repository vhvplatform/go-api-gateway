package circuitbreaker

import (
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
