package middleware

import (
"bytes"
"net/http"
"net/http/httptest"
"testing"

"github.com/gin-gonic/gin"
"github.com/stretchr/testify/assert"
)

func TestRequestValidationMiddleware(t *testing.T) {
gin.SetMode(gin.TestMode)

t.Run("allows GET requests without Content-Type", func(t *testing.T) {
r := gin.New()
r.Use(RequestValidationMiddleware())
r.GET("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("GET", "/test", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("allows DELETE requests without Content-Type", func(t *testing.T) {
r := gin.New()
r.Use(RequestValidationMiddleware())
r.DELETE("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("DELETE", "/test", nil)
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
})

t.Run("rejects POST without Content-Type", func(t *testing.T) {
r := gin.New()
r.Use(RequestValidationMiddleware())
r.POST("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("{}"))
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

assert.Equal(t, http.StatusBadRequest, w.Code)
})

t.Run("allows POST with Content-Type", func(t *testing.T) {
r := gin.New()
r.Use(RequestValidationMiddleware())
r.POST("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

req := httptest.NewRequest("POST", "/test", bytes.NewBufferString("{}"))
req.Header.Set("Content-Type", "application/json")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
})
}

func TestRequestSizeLimitMiddleware(t *testing.T) {
gin.SetMode(gin.TestMode)

t.Run("allows request within size limit", func(t *testing.T) {
r := gin.New()
r.Use(RequestSizeLimitMiddleware(1024)) // 1KB limit
r.POST("/test", func(c *gin.Context) {
c.JSON(http.StatusOK, gin.H{"message": "ok"})
})

smallBody := bytes.NewBufferString("small body")
req := httptest.NewRequest("POST", "/test", smallBody)
req.Header.Set("Content-Type", "application/json")
w := httptest.NewRecorder()
r.ServeHTTP(w, req)

assert.Equal(t, http.StatusOK, w.Code)
})
}
