package middleware

import (
"net/http"

"github.com/gin-gonic/gin"
)

// RequestSizeLimitMiddleware limits the size of incoming requests
func RequestSizeLimitMiddleware(maxBytes int64) gin.HandlerFunc {
return func(c *gin.Context) {
c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBytes)
c.Next()
}
}

// RequestValidationMiddleware validates incoming requests
func RequestValidationMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
// Validate required headers for non-GET/DELETE requests
if c.Request.Method != http.MethodGet && c.Request.Method != http.MethodDelete {
contentType := c.GetHeader("Content-Type")
if contentType == "" {
c.JSON(http.StatusBadRequest, gin.H{
"error": "Content-Type header required",
"code":  "MISSING_CONTENT_TYPE",
})
c.Abort()
return
}
}
c.Next()
}
}
