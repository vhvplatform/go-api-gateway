package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-api-gateway/internal/errors"
)

// TimeoutMiddleware adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()

		// Check if timeout occurred
		if ctx.Err() == context.DeadlineExceeded {
			// Only send error response if nothing was written yet
			if !c.Writer.Written() {
				correlationID := c.GetString("correlation_id")
				errorResp := errors.NewErrorResponse(
					"TIMEOUT",
					"Request timeout exceeded",
					nil,
					correlationID,
				)
				c.JSON(http.StatusGatewayTimeout, errorResp)
			}
			c.Abort()
		}
	}
}
