package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-api-gateway/internal/middleware"
	"github.com/vhvplatform/go-shared/logger"
)

// SetupPermissionExampleRoutes demonstrates how to use permission middleware
// This is an example of how to protect routes with specific permissions
func SetupPermissionExampleRoutes(
	r *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	permMiddleware *middleware.PermissionMiddleware,
	log *logger.Logger,
) {
	// Example: User Management Routes with Permissions
	users := r.Group("/api/users")
	users.Use(authMiddleware.Authenticate()) // First verify token
	{
		// GET /api/users - Requires "user.read" permission
		users.GET("",
			permMiddleware.RequirePermission("user.read"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "List all users",
					"data":    []string{"user1", "user2", "user3"},
				})
			})

		// GET /api/users/:id - Requires "user.read" permission
		users.GET("/:id",
			permMiddleware.RequirePermission("user.read"),
			func(c *gin.Context) {
				userID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"message": "Get user details",
					"user_id": userID,
				})
			})

		// POST /api/users - Requires BOTH "user.write" AND "user.create"
		users.POST("",
			permMiddleware.RequirePermission("user.write", "user.create"),
			func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{
					"message": "User created successfully",
					"user_id": "new-user-123",
				})
			})

		// PUT /api/users/:id - Requires "user.write" permission
		users.PUT("/:id",
			permMiddleware.RequirePermission("user.write"),
			func(c *gin.Context) {
				userID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"message": "User updated successfully",
					"user_id": userID,
				})
			})

		// DELETE /api/users/:id - Requires "user.delete" permission
		users.DELETE("/:id",
			permMiddleware.RequirePermission("user.delete"),
			func(c *gin.Context) {
				userID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"message": "User deleted successfully",
					"user_id": userID,
				})
			})
	}

	// Example: Tenant Management Routes with Permissions
	tenants := r.Group("/api/tenants")
	tenants.Use(authMiddleware.Authenticate())
	{
		// GET /api/tenants - Requires "tenant.read"
		tenants.GET("",
			permMiddleware.RequirePermission("tenant.read"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "List all tenants",
					"data":    []string{"tenant1", "tenant2"},
				})
			})

		// POST /api/tenants - Requires "tenant.manage"
		tenants.POST("",
			permMiddleware.RequirePermission("tenant.manage"),
			func(c *gin.Context) {
				c.JSON(http.StatusCreated, gin.H{
					"message": "Tenant created successfully",
				})
			})

		// PUT /api/tenants/:id/settings - Requires "tenant.manage"
		tenants.PUT("/:id/settings",
			permMiddleware.RequirePermission("tenant.manage"),
			func(c *gin.Context) {
				tenantID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"message":   "Tenant settings updated",
					"tenant_id": tenantID,
				})
			})

		// DELETE /api/tenants/:id - Requires "tenant.delete"
		tenants.DELETE("/:id",
			permMiddleware.RequirePermission("tenant.delete"),
			func(c *gin.Context) {
				tenantID := c.Param("id")
				c.JSON(http.StatusOK, gin.H{
					"message":   "Tenant deleted successfully",
					"tenant_id": tenantID,
				})
			})
	}

	// Example: Admin Routes - Any of admin permissions
	admin := r.Group("/api/admin")
	admin.Use(authMiddleware.Authenticate())
	{
		// Dashboard - Requires admin OR super_admin permission
		admin.GET("/dashboard",
			permMiddleware.RequireAnyPermission("admin.dashboard", "super_admin.*"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Admin dashboard",
					"stats": gin.H{
						"users":   100,
						"tenants": 10,
					},
				})
			})

		// System Config - Requires super_admin role
		admin.GET("/system/config",
			permMiddleware.RequireRole("super_admin"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "System configuration",
					"config": gin.H{
						"maintenance_mode": false,
						"api_version":      "v1",
					},
				})
			})

		// Audit Logs - Requires "system.audit" permission
		admin.GET("/audit-logs",
			permMiddleware.RequirePermission("system.audit"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Audit logs",
					"logs":    []string{"log1", "log2"},
				})
			})
	}

	// Example: Public Routes (no authentication/permission required)
	public := r.Group("/api/public")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "healthy",
				"time":   "2026-01-08T00:00:00Z",
			})
		})

		public.GET("/version", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"version": "1.0.0",
				"build":   "20260108",
			})
		})
	}

	// Example: Optional Auth Routes (authenticated users get more data)
	optional := r.Group("/api/content")
	optional.Use(authMiddleware.OptionalAuthenticate())
	{
		optional.GET("/articles", func(c *gin.Context) {
			userID, authenticated := middleware.GetUserID(c)

			response := gin.H{
				"message": "List articles",
				"data":    []string{"article1", "article2"},
			}

			if authenticated {
				response["user_id"] = userID
				response["premium_content"] = true
			}

			c.JSON(http.StatusOK, response)
		})
	}

	log.Info("Permission example routes configured successfully")
}

// SetupPermissionTestRoutes creates test endpoints for permission system
func SetupPermissionTestRoutes(
	r *gin.Engine,
	authMiddleware *middleware.AuthMiddleware,
	permMiddleware *middleware.PermissionMiddleware,
) {
	test := r.Group("/api/test")
	test.Use(authMiddleware.Authenticate())
	{
		// Test wildcard permission matching
		test.GET("/wildcard",
			permMiddleware.RequirePermission("test.*"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Wildcard permission test passed",
				})
			})

		// Test multiple permissions
		test.POST("/multiple",
			permMiddleware.RequirePermission("test.read", "test.write", "test.execute"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Multiple permission test passed",
				})
			})

		// Test any permission
		test.GET("/any",
			permMiddleware.RequireAnyPermission("test.option1", "test.option2", "test.option3"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Any permission test passed",
				})
			})

		// Test role-based access
		test.GET("/role",
			permMiddleware.RequireRole("tester", "developer"),
			func(c *gin.Context) {
				c.JSON(http.StatusOK, gin.H{
					"message": "Role test passed",
				})
			})
	}
}

// RoutePermissionMap maps route patterns to required permissions
// This can be used for automatic permission-based routing
var RoutePermissionMap = map[string][]string{
	"GET:/api/users":        {"user.read"},
	"POST:/api/users":       {"user.write", "user.create"},
	"PUT:/api/users/:id":    {"user.write"},
	"DELETE:/api/users/:id": {"user.delete"},

	"GET:/api/tenants":        {"tenant.read"},
	"POST:/api/tenants":       {"tenant.manage"},
	"PUT:/api/tenants/:id":    {"tenant.manage"},
	"DELETE:/api/tenants/:id": {"tenant.delete"},

	"GET:/api/admin/dashboard":     {"admin.dashboard", "super_admin.*"},
	"GET:/api/admin/system/config": {"system.config"},
	"GET:/api/admin/audit-logs":    {"system.audit"},
}
