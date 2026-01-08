package router

import (
	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-api-gateway/internal/client"
	"github.com/vhvplatform/go-api-gateway/internal/handler"
	internalmiddleware "github.com/vhvplatform/go-api-gateway/internal/middleware"
	"github.com/vhvplatform/go-shared/cache"
	"github.com/vhvplatform/go-shared/config"
	"github.com/vhvplatform/go-shared/logger"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authClient *client.AuthClient,
	cacheClient cache.Cache,
	proxyHandler *handler.ProxyHandler,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	tenantHandler *handler.TenantHandler,
	notificationHandler *handler.NotificationHandler,
	log *logger.Logger,
) {
	// 1. PUBLIC API ROUTES
	public := r.Group("/auth")
	{
		public.POST("/login", authHandler.Login)
		public.POST("/register", authHandler.Register)
		public.POST("/refresh", authHandler.RefreshToken)
	}

	// 2. PROTECTED DYNAMIC API ROUTES (/api/:service/*path)
	api := r.Group("/api")
	api.Use(internalmiddleware.AuthMiddleware(authClient, cacheClient, cfg.JWT.Secret))
	{
		// This handles /api/user/profile, /api/tenant/settings, etc.
		// The "service" param is used by proxyHandler.APIProxy
		api.Any("/:service/*path", proxyHandler.APIProxy)
	}

	// 3. PAGE ROUTES (/page/service-name/page-path)
	r.GET("/page/*path", proxyHandler.PageProxy)

	// 4. UPLOAD ROUTES (/upload/file-key)
	r.Any("/upload/*path", proxyHandler.UploadProxy)

	// 5. FALLBACK / BEAUTIFUL URLS (Slug)
	r.NoRoute(proxyHandler.SlugProxy)

	log.Info("Universal dynamic routes configured successfully")
}
