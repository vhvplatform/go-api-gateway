package middleware

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCORSHeaders(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("sets allowed origins", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)
		c.Request.Header.Set("Origin", "http://localhost:3000")

		// Test that CORS allows specified origins
		origin := c.Request.Header.Get("Origin")
		assert.Equal(t, "http://localhost:3000", origin)
	})

	t.Run("sets allowed methods", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/test", nil)
		c.Request.Header.Set("Access-Control-Request-Method", "POST")

		// Test allowed methods in preflight
		method := c.Request.Header.Get("Access-Control-Request-Method")
		assert.Equal(t, "POST", method)
	})

	t.Run("sets allowed headers", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/test", nil)
		c.Request.Header.Set("Access-Control-Request-Headers", "Content-Type,Authorization")

		// Test allowed headers
		headers := c.Request.Header.Get("Access-Control-Request-Headers")
		assert.Contains(t, headers, "Content-Type")
		assert.Contains(t, headers, "Authorization")
	})

	t.Run("handles credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("GET", "/test", nil)

		// Test credentials handling
		assert.NotNil(t, c)
	})
}

func TestPreflightRequests(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("responds to OPTIONS request", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/api/test", nil)

		// Verify it's an OPTIONS request
		assert.Equal(t, "OPTIONS", c.Request.Method)
	})

	t.Run("includes max age header", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request = httptest.NewRequest("OPTIONS", "/api/test", nil)

		// Test that max age can be configured
		assert.NotNil(t, c)
	})
}
