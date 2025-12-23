package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// RateLimiter implements rate limiting
type RateLimiter struct {
	limiters map[string]*rate.Limiter
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*rate.Limiter),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// GetLimiter returns a limiter for the given key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	limiter, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if !exists {
		rl.mu.Lock()
		limiter = rate.NewLimiter(rl.rate, rl.burst)
		rl.limiters[key] = limiter
		rl.mu.Unlock()
	}

	return limiter
}

// CleanupLimiters removes old limiters
func (rl *RateLimiter) CleanupLimiters() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for key := range rl.limiters {
			delete(rl.limiters, key)
		}
		rl.mu.Unlock()
	}
}

// RateLimitMiddleware implements rate limiting middleware
func RateLimitMiddleware(rl *RateLimiter) gin.HandlerFunc {
	// Start cleanup goroutine
	go rl.CleanupLimiters()

	return func(c *gin.Context) {
		// Use IP address as the key
		key := c.ClientIP()

		// Get tenant ID if available for per-tenant limiting
		tenantID := c.GetHeader("X-Tenant-ID")
		if tenantID != "" {
			key = tenantID + ":" + key
		}

		limiter := rl.GetLimiter(key)

		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
