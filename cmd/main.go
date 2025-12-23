package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/config"
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/client"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/handler"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/middleware"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/router"
)

func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	log.Info("Starting API Gateway...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", "error", err)
	}

	// Initialize gRPC clients
	authClient := client.NewAuthClient(getServiceURL("AUTH_SERVICE_URL", "auth-service:50051"), log)
	userClient := client.NewUserClient(getServiceURL("USER_SERVICE_URL", "user-service:50052"), log)
	tenantClient := client.NewTenantClient(getServiceURL("TENANT_SERVICE_URL", "tenant-service:50053"), log)

	// Initialize HTTP client for notification service
	notificationURL := getServiceURL("NOTIFICATION_SERVICE_URL", "http://notification-service:8084")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authClient, log)
	userHandler := handler.NewUserHandler(userClient, log)
	tenantHandler := handler.NewTenantHandler(tenantClient, log)
	notificationHandler := handler.NewNotificationHandler(notificationURL, log)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.LoggerMiddleware(log))
	r.Use(middleware.CorrelationIDMiddleware())

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Correlation-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Correlation-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Rate limiting middleware
	rateLimiter := middleware.NewRateLimiter(100, 200) // 100 req/s, burst 200
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Health check endpoints
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Setup routes
	router.SetupRoutes(r, cfg, authHandler, userHandler, tenantHandler, notificationHandler, log)

	// Start HTTP server
	port := os.Getenv("API_GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Info("API Gateway started", "port", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", "error", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API Gateway...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	// Close gRPC connections
	authClient.Close()
	userClient.Close()
	tenantClient.Close()

	log.Info("API Gateway stopped")
}

func getServiceURL(envVar, defaultValue string) string {
	url := os.Getenv(envVar)
	if url == "" {
		return defaultValue
	}
	return url
}
