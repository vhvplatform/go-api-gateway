package middleware

import (
"net/http"

"github.com/gin-gonic/gin"
"github.com/vhvcorp/go-shared/logger"
"github.com/vhvcorp/go-api-gateway/internal/errors"
"go.uber.org/zap"
)

// RecoveryMiddleware provides panic recovery with proper logging
func RecoveryMiddleware(log *logger.Logger) gin.HandlerFunc {
return func(c *gin.Context) {
defer func() {
if err := recover(); err != nil {
log.Error("Panic recovered", zap.Any("error", err))

correlationID := c.GetString("correlation_id")
errorResp := errors.NewErrorResponse(
"INTERNAL_ERROR",
"An internal server error occurred",
nil,
correlationID,
)

c.JSON(http.StatusInternalServerError, errorResp)
c.Abort()
}
}()
c.Next()
}
}
