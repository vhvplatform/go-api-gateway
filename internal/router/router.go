package router

import (
	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/config"
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/handler"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	notificationHandler *handler.NotificationHandler,
	log *logger.Logger,
) {
	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", middleware.AuthMiddleware(cfg.JWT.Secret), authHandler.Logout)
		}

		// Protected routes (require authentication)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg.JWT.Secret))
		{
			// User routes
			users := protected.Group("/users")
			{
				users.POST("", userHandler.CreateUser)
				users.GET("", userHandler.GetUsers)
				users.GET("/search", userHandler.SearchUsers)
				users.GET("/:id", userHandler.GetUser)
				users.PUT("/:id", userHandler.UpdateUser)
				users.DELETE("/:id", userHandler.DeleteUser)
			}

			// Tenant routes
			tenants := protected.Group("/tenants")
			{
				tenants.POST("", tenantHandler.CreateTenant)
				tenants.GET("", tenantHandler.GetTenants)
				tenants.GET("/:id", tenantHandler.GetTenant)
				tenants.PUT("/:id", tenantHandler.UpdateTenant)
				tenants.DELETE("/:id", tenantHandler.DeleteTenant)
				tenants.POST("/:id/users", tenantHandler.AddUserToTenant)
				tenants.DELETE("/:id/users/:user_id", tenantHandler.RemoveUserFromTenant)
			}

			// Notification routes
			notifications := protected.Group("/notifications")
			{
				notifications.POST("/email", notificationHandler.SendEmail)
				notifications.POST("/webhook", notificationHandler.SendWebhook)
				notifications.GET("", notificationHandler.GetNotifications)
				notifications.GET("/:id", notificationHandler.GetNotification)
			}
		}
	}

	log.Info("Routes configured successfully")
}
