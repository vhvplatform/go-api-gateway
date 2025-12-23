package middleware

import (
"context"
"net/http"
"time"

"github.com/gin-gonic/gin"
"github.com/longvhv/saas-framework-go/services/api-gateway/internal/errors"
)

// TimeoutMiddleware adds a timeout to requests
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
return func(c *gin.Context) {
ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
defer cancel()

c.Request = c.Request.WithContext(ctx)

done := make(chan struct{})
go func() {
c.Next()
close(done)
}()

select {
case <-done:
return
case <-ctx.Done():
if ctx.Err() == context.DeadlineExceeded {
correlationID := c.GetString("correlation_id")
errorResp := errors.NewErrorResponse(
"TIMEOUT",
"Request timeout exceeded",
nil,
correlationID,
)
c.JSON(http.StatusGatewayTimeout, errorResp)
c.Abort()
}
}
}
}
