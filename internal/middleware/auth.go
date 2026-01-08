package middleware

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-api-gateway/internal/client"
	"github.com/vhvplatform/go-shared/cache"
	"github.com/vhvplatform/go-shared/jwt"
)

// AuthMiddleware validates Opaque tokens via AuthService and injects Internal JWT
func AuthMiddleware(authClient *client.AuthClient, tieredCache cache.Cache, jwtSecret string) gin.HandlerFunc {
	jwtManager := jwt.NewManager(jwtSecret, 3600, 86400)

	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		opaqueToken := parts[1]
		cacheKey := "token:" + opaqueToken

		// 1. Check Cache (Level 1 & Level 2 handled by TieredCache)
		var resp client.VerifyTokenResponse
		err := tieredCache.Get(c.Request.Context(), cacheKey, &resp)
		if err == nil && resp.Valid {
			// Cache hit
			injectHeaders(c, &resp, jwtManager)
			c.Next()
			return
		}

		// 2. Cache miss, Verify Opaque Token via gRPC
		apiResp, err := authClient.VerifyToken(c.Request.Context(), opaqueToken)
		if err != nil || !apiResp.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// 3. Cache the result (TieredCache handles L1/L2 updates)
		_ = tieredCache.Set(c.Request.Context(), cacheKey, apiResp, 30*time.Minute)

		// 4. Inject Headers and proceed
		injectHeaders(c, apiResp, jwtManager)
		c.Next()
	}
}

func injectHeaders(c *gin.Context, resp *client.VerifyTokenResponse, jwtManager *jwt.Manager) {
	// Generate Internal JWT
	// Signature: GenerateToken(userID, tenantID, email string, roles, permissions []string)
	internalToken, _ := jwtManager.GenerateToken(resp.UserId, resp.TenantId, resp.Email, resp.Roles, resp.Permissions)

	// Inject Headers for backend services
	c.Request.Header.Set("X-Tenant-ID", resp.TenantId)
	c.Request.Header.Set("X-Internal-Token", internalToken)

	// Set context vars for gateway internal use
	c.Set("user_id", resp.UserId)
	c.Set("tenant_id", resp.TenantId)
	c.Set("roles", resp.Roles)
	c.Set("permissions", resp.Permissions)
}

// TenantMiddleware ensures tenant context is available
func TenantMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID := c.GetString("tenant_id")
		if tenantID == "" {
			// Try to get from header
			tenantID = c.GetHeader("X-Tenant-ID")
			if tenantID == "" {
				c.JSON(http.StatusBadRequest, gin.H{
					"error": "Tenant ID required",
				})
				c.Abort()
				return
			}
			c.Set("tenant_id", tenantID)
		}

		c.Next()
	}
}
