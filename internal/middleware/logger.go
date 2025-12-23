package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/logger"
)

// LoggerMiddleware logs HTTP requests
func LoggerMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		log.Info("HTTP Request",
			"method", c.Request.Method,
			"path", path,
			"query", query,
			"status", c.Writer.Status(),
			"latency", latency.String(),
			"ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"correlation_id", c.GetString("correlation_id"),
		)
	}
}
