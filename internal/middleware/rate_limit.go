package middleware

import (
	"context"
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
	_, exists := rl.limiters[key]
	rl.mu.RUnlock()

	if exists {
		// Update last access time
		rl.mu.Lock()
		// Check again after acquiring write lock
		if entry, exists := rl.limiters[key]; exists {
			entry.lastAccess = time.Now()
			rl.mu.Unlock()
			return entry.limiter
		}
		rl.mu.Unlock()
	}

	rl.mu.Lock()
	// Double-check after acquiring write lock
	if entry, exists := rl.limiters[key]; exists {
		entry.lastAccess = time.Now()
		rl.mu.Unlock()
		return entry.limiter
	}

	limiter := rate.NewLimiter(rl.rate, rl.burst)
	rl.limiters[key] = &limiterEntry{
		limiter:    limiter,
		lastAccess: time.Now(),
	}
	rl.mu.Unlock()

	return limiter
}

// CleanupLimiters removes inactive limiters
func (rl *RateLimiter) CleanupLimiters(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			rl.mu.Lock()
			now := time.Now()
			cleaned := 0
			for key, entry := range rl.limiters {
				// Only delete if inactive for 10 minutes
				if now.Sub(entry.lastAccess) > 10*time.Minute {
					delete(rl.limiters, key)
					cleaned++
				}
			}
			rl.mu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

// RateLimitMiddleware implements rate limiting middleware
func RateLimitMiddleware(rl *RateLimiter, ctx context.Context) gin.HandlerFunc {
	// Start cleanup goroutine with context
	go rl.CleanupLimiters(ctx)

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
