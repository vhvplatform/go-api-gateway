package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/auth"
	"github.com/vhvplatform/go-shared/cache"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
)

// PermissionConfig holds configuration for permission middleware
type PermissionConfig struct {
	// AuthClient is the gRPC client for auth service
	AuthClient interface {
		CheckPermission(ctx context.Context, userID, tenantID, permission string) (bool, error)
		GetUserRoles(ctx context.Context, userID, tenantID string) ([]string, error)
	}
	// Cache is the 2-level cache (L1 local + L2 Redis)
	Cache cache.Cache
	// Logger for logging
	Logger *logger.Logger
	// CacheTTL is how long to cache permissions (default: 5 minutes)
	CacheTTL time.Duration
	// SkipPaths are paths that don't require permission checks
	SkipPaths []string
}

// PermissionMiddleware creates a middleware that checks user permissions
type PermissionMiddleware struct {
	config *PermissionConfig
}

// NewPermissionMiddleware creates a new permission middleware
func NewPermissionMiddleware(config *PermissionConfig) *PermissionMiddleware {
	if config.CacheTTL == 0 {
		config.CacheTTL = 5 * time.Minute
	}

	return &PermissionMiddleware{
		config: config,
	}
}

// RequirePermission creates a middleware that requires specific permissions
func (m *PermissionMiddleware) RequirePermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		if m.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get user context from previous auth middleware
		userID, exists := c.Get("user_id")
		if !exists {
			m.config.Logger.Warn("Permission check failed: no user_id in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		tenantID, exists := c.Get("tenant_id")
		if !exists {
			m.config.Logger.Warn("Permission check failed: no tenant_id in context")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
			c.Abort()
			return
		}

		userIDStr := userID.(string)
		tenantIDStr := tenantID.(string)

		// Check all required permissions
		hasPermission, missing, err := m.checkPermissions(c.Request.Context(), userIDStr, tenantIDStr, permissions)
		if err != nil {
			m.config.Logger.Error("Permission check error",
				zap.String("user_id", userIDStr),
				zap.String("tenant_id", tenantIDStr),
				zap.Strings("required_permissions", permissions),
				zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "permission check failed"})
			c.Abort()
			return
		}

		if !hasPermission {
			m.config.Logger.Warn("Permission denied",
				zap.String("user_id", userIDStr),
				zap.String("tenant_id", tenantIDStr),
				zap.Strings("required_permissions", permissions),
				zap.Strings("missing_permissions", missing))
			c.JSON(http.StatusForbidden, gin.H{
				"error":                "insufficient permissions",
				"required_permissions": permissions,
				"missing_permissions":  missing,
			})
			c.Abort()
			return
		}

		// Permission granted, continue
		m.config.Logger.Debug("Permission granted",
			zap.String("user_id", userIDStr),
			zap.String("tenant_id", tenantIDStr),
			zap.Strings("permissions", permissions))
		c.Next()
	}
}

// RequireAnyPermission creates a middleware that requires at least one of the permissions
func (m *PermissionMiddleware) RequireAnyPermission(permissions ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if path should be skipped
		if m.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		// Get user context
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			c.Abort()
			return
		}

		tenantID, exists := c.Get("tenant_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "tenant context required"})
			c.Abort()
			return
		}

		userIDStr := userID.(string)
		tenantIDStr := tenantID.(string)

		// Check if user has any of the required permissions
		hasAny, err := m.checkAnyPermission(c.Request.Context(), userIDStr, tenantIDStr, permissions)
		if err != nil {
			m.config.Logger.Error("Permission check error", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "permission check failed"})
			c.Abort()
			return
		}

		if !hasAny {
			m.config.Logger.Warn("Permission denied - none of required permissions",
				zap.String("user_id", userIDStr),
				zap.String("tenant_id", tenantIDStr),
				zap.Strings("any_of_permissions", permissions))
			c.JSON(http.StatusForbidden, gin.H{
				"error":  "insufficient permissions",
				"any_of": permissions,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRole creates a middleware that requires a specific role
func (m *PermissionMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if m.shouldSkipPath(c.Request.URL.Path) {
			c.Next()
			return
		}

		userID, _ := c.Get("user_id")
		tenantID, _ := c.Get("tenant_id")
		userIDStr := userID.(string)
		tenantIDStr := tenantID.(string)

		// Get user roles from cache or auth service
		userRoles, err := m.getUserRoles(c.Request.Context(), userIDStr, tenantIDStr)
		if err != nil {
			m.config.Logger.Error("Failed to get user roles", zap.Error(err))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "role check failed"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, requiredRole := range roles {
			for _, userRole := range userRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			m.config.Logger.Warn("Role check failed",
				zap.String("user_id", userIDStr),
				zap.String("tenant_id", tenantIDStr),
				zap.Strings("required_roles", roles),
				zap.Strings("user_roles", userRoles))
			c.JSON(http.StatusForbidden, gin.H{
				"error":          "insufficient role",
				"required_roles": roles,
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// PermissionFromRoute extracts permission from route metadata
// Routes should define permissions like:
//
//	router.GET("/users", middleware.PermissionFromRoute(), handler)
//
// And route metadata should include: "permission": "user.read"
func (m *PermissionMiddleware) PermissionFromRoute() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get permission from route metadata
		permission, exists := c.Get("route_permission")
		if !exists || permission == "" {
			// No permission required for this route
			c.Next()
			return
		}

		permStr := permission.(string)

		// Use RequirePermission to check
		m.RequirePermission(permStr)(c)
	}
}

// Helper methods

func (m *PermissionMiddleware) shouldSkipPath(path string) bool {
	for _, skipPath := range m.config.SkipPaths {
		if strings.HasPrefix(path, skipPath) {
			return true
		}
	}
	return false
}

func (m *PermissionMiddleware) checkPermissions(ctx context.Context, userID, tenantID string, permissions []string) (bool, []string, error) {
	// Get all user permissions from cache
	userPermissions, err := m.getUserPermissions(ctx, userID, tenantID)
	if err != nil {
		return false, nil, err
	}

	// Create permission set
	permSet, err := auth.NewPermissionSet(userPermissions)
	if err != nil {
		return false, nil, err
	}

	// Check each required permission
	missing := []string{}
	for _, required := range permissions {
		if !permSet.Has(required) {
			missing = append(missing, required)
		}
	}

	hasAll := len(missing) == 0
	return hasAll, missing, nil
}

func (m *PermissionMiddleware) checkAnyPermission(ctx context.Context, userID, tenantID string, permissions []string) (bool, error) {
	userPermissions, err := m.getUserPermissions(ctx, userID, tenantID)
	if err != nil {
		return false, err
	}

	permSet, err := auth.NewPermissionSet(userPermissions)
	if err != nil {
		return false, err
	}

	return permSet.HasAny(permissions...), nil
}

func (m *PermissionMiddleware) getUserPermissions(ctx context.Context, userID, tenantID string) ([]string, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("permissions:%s:%s", userID, tenantID)
	var cachedPerms []string

	if m.config.Cache != nil {
		err := m.config.Cache.Get(ctx, cacheKey, &cachedPerms)
		if err == nil && len(cachedPerms) > 0 {
			m.config.Logger.Debug("Permission cache hit",
				zap.String("user_id", userID),
				zap.String("tenant_id", tenantID))
			return cachedPerms, nil
		}
	}

	// Cache miss - call auth service via gRPC
	m.config.Logger.Debug("Permission cache miss, calling auth service",
		zap.String("user_id", userID),
		zap.String("tenant_id", tenantID))

	// Note: In real implementation, we would need to query all permissions
	// For now, we'll return empty and rely on CheckPermission calls
	// A better approach would be to add a GetUserPermissions gRPC method
	permissions := []string{}

	// Cache the result
	if m.config.Cache != nil {
		_ = m.config.Cache.Set(ctx, cacheKey, permissions, m.config.CacheTTL)
	}

	return permissions, nil
}

func (m *PermissionMiddleware) getUserRoles(ctx context.Context, userID, tenantID string) ([]string, error) {
	// Try cache first
	cacheKey := fmt.Sprintf("roles:%s:%s", userID, tenantID)
	var cachedRoles []string

	if m.config.Cache != nil {
		err := m.config.Cache.Get(ctx, cacheKey, &cachedRoles)
		if err == nil && len(cachedRoles) > 0 {
			return cachedRoles, nil
		}
	}

	// Cache miss - call auth service
	if m.config.AuthClient != nil {
		roles, err := m.config.AuthClient.GetUserRoles(ctx, userID, tenantID)
		if err != nil {
			return nil, fmt.Errorf("failed to get user roles: %w", err)
		}

		// Cache the result
		if m.config.Cache != nil {
			_ = m.config.Cache.Set(ctx, cacheKey, roles, m.config.CacheTTL)
		}

		return roles, nil
	}

	return []string{}, nil
}
