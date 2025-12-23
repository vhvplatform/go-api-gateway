package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// limiterEntry holds a rate limiter and its last access time
type limiterEntry struct {
	limiter    *rate.Limiter
	lastAccess time.Time
}

// RateLimiter implements rate limiting
type RateLimiter struct {
	limiters map[string]*limiterEntry
	mu       sync.RWMutex
	rate     rate.Limit
	burst    int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(rps float64, burst int) *RateLimiter {
	return &RateLimiter{
		limiters: make(map[string]*limiterEntry),
		rate:     rate.Limit(rps),
		burst:    burst,
	}
}

// GetLimiter returns a limiter for the given key
func (rl *RateLimiter) GetLimiter(key string) *rate.Limiter {
	rl.mu.RLock()
	entry, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if exists {
		// Update last access time
		rl.mu.Lock()
		entry.lastAccess = time.Now()
		rl.mu.Unlock()
		return entry.limiter
	}

	rl.mu.Lock()
	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[key] = &limiterEntry{
		limiter:    limiter,
		lastAccess: time.Now(),
	}
	rl.mu.Unlock()

	return limiter
}

// CleanupLimiters removes inactive limiters
func (rl *RateLimiter) CleanupLimiters() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, entry := range rl.limiters {
			// Only delete if inactive for 10 minutes
			if now.Sub(entry.lastAccess) > 10*time.Minute {
				delete(rl.limiters, key)
			}
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
