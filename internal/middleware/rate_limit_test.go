package middleware

import (
"net/http"
"net/http/httptest"
"testing"
"time"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
)

func TestRateLimiter(t *testing.T) {
gin.SetMode(gin.TestMode)

t.Run("allows requests within limit", func(t *testing.T) {
rl := NewRateLimiter(10, 10)
limiter := rl.GetLimiter("test-key")

// First 10 requests should be allowed
for i := 0; i < 10; i++ {
assert.True(t, limiter.Allow(), "Request %d should be allowed", i+1)
}
})

t.Run("blocks requests exceeding limit", func(t *testing.T) {
rl := NewRateLimiter(1, 1)
limiter := rl.GetLimiter("test-key-2")

// First request allowed
assert.True(t, limiter.Allow())

// Second request blocked (burst used up)
assert.False(t, limiter.Allow())
})

t.Run("tracks last access time", func(t *testing.T) {
rl := NewRateLimiter(10, 10)

limiter1 := rl.GetLimiter("test-key-3")
assert.NotNil(t, limiter1)

time.Sleep(10 * time.Millisecond)

limiter2 := rl.GetLimiter("test-key-3")
assert.Equal(t, limiter1, limiter2, "Should return same limiter for same key")
})
}

func TestRateLimitMiddleware(t *testing.T) {
gin.SetMode(gin.TestMode)

t.Run("allows requests within limit", func(t *testing.T) {
rl := NewRateLimiter(100, 100)

r := gin.New()
r.Use(RateLimitMiddleware(rl))
r.GET("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("GET", "/test", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("blocks requests exceeding limit", func(t *testing.T) {
rl := NewRateLimiter(1, 1)

r := gin.New()
r.Use(RateLimitMiddleware(rl))
r.GET("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

// First request should succeed
req1 := httptest.NewRequest("GET", "/test", nil)
req1.RemoteAddr = "192.168.1.1:1234"
w1 := httptest.NewRecorder()
r.ServeHTTP(w1, req1)
assert.Equal(t, http.StatusOK, w1.Code)

// Second request should be rate limited
req2 := httptest.NewRequest("GET", "/test", nil)
req2.RemoteAddr = "192.168.1.1:1234"
w2 := httptest.NewRecorder()
r.ServeHTTP(w2, req2)
assert.Equal(t, http.StatusTooManyRequests, w2.Code)
})
}
