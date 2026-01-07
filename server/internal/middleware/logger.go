package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
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
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.Int("status", c.Writer.Status()),
			zap.String("latency", latency.String()),
			zap.String("ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
			zap.String("correlation_id", c.GetString("correlation_id")),
		)
	}
}
