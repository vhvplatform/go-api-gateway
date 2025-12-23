package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/config"
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/cache"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/circuitbreaker"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/client"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/handler"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/health"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/middleware"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/router"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/tracing"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

	// Initialize distributed tracing (optional)
	if os.Getenv("ENABLE_TRACING") == "true" {
		jaegerURL := getServiceURL("JAEGER_URL", "http://jaeger:14268/api/traces")
		tp, err := tracing.InitTracer("api-gateway", jaegerURL)
		if err != nil {
			log.Error("Failed to initialize tracer", "error", err)
		} else {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := tp.Shutdown(ctx); err != nil {
					log.Error("Failed to shutdown tracer", "error", err)
				}
			}()
			log.Info("Distributed tracing enabled")
		}
	}

	// Initialize Redis cache (optional)
	var cacheClient *cache.Cache
	if redisURL := os.Getenv("REDIS_URL"); redisURL != "" {
		cacheClient, err = cache.NewCache(redisURL)
		if err != nil {
			log.Error("Failed to initialize cache", "error", err)
		} else {
			defer cacheClient.Close()
			log.Info("Redis cache initialized")
		}
	}

	// Initialize circuit breaker
	cb := circuitbreaker.NewCircuitBreaker()

	// Initialize health checker
	healthChecker := health.NewHealthChecker()
	// Register health checks for services
	healthChecker.RegisterCheck("auth-service", func(ctx context.Context) error {
		// TODO: Implement actual health check when proto is available
		return nil
	})
	healthChecker.RegisterCheck("user-service", func(ctx context.Context) error {
		// TODO: Implement actual health check when proto is available
		return nil
	})
	healthChecker.RegisterCheck("tenant-service", func(ctx context.Context) error {
		// TODO: Implement actual health check when proto is available
		return nil
	})

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

	// Recovery middleware with custom error handling
	r.Use(middleware.RecoveryMiddleware(log))

	// Correlation ID middleware (should be first)
	r.Use(middleware.CorrelationIDMiddleware())

	// Logging middleware
	r.Use(middleware.LoggerMiddleware(log))

	// Metrics middleware (if enabled)
	if os.Getenv("ENABLE_METRICS") != "false" {
		r.Use(middleware.MetricsMiddleware())
	}

	// Compression middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Request validation middleware
	r.Use(middleware.RequestValidationMiddleware())

	// Request size limit middleware
	maxRequestSize := int64(10485760) // 10MB default
	if size := os.Getenv("MAX_REQUEST_SIZE"); size != "" {
		if parsedSize, err := strconv.ParseInt(size, 10, 64); err == nil {
			maxRequestSize = parsedSize
		}
	}
	r.Use(middleware.RequestSizeLimitMiddleware(maxRequestSize))

	// Timeout middleware
	r.Use(middleware.TimeoutMiddleware(30 * time.Second))

	// CORS configuration
	corsConfig := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Correlation-ID", "X-Tenant-ID"},
		ExposeHeaders:    []string{"Content-Length", "X-Correlation-ID"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(corsConfig))

	// Rate limiting middleware
	rateLimit := 100.0
	rateBurst := 200
	if rps := os.Getenv("RATE_LIMIT_RPS"); rps != "" {
		if parsedRPS, err := strconv.ParseFloat(rps, 64); err == nil {
			rateLimit = parsedRPS
		}
	}
	if burst := os.Getenv("RATE_LIMIT_BURST"); burst != "" {
		if parsedBurst, err := strconv.Atoi(burst); err == nil {
			rateBurst = parsedBurst
		}
	}
	rateLimiter := middleware.NewRateLimiter(rateLimit, rateBurst)
	r.Use(middleware.RateLimitMiddleware(rateLimiter))

	// Health check endpoints
	r.GET("/health", func(c *gin.Context) {
		status := healthChecker.CheckAll(c.Request.Context())
		if status.Status == "healthy" {
			c.JSON(http.StatusOK, status)
		} else {
			c.JSON(http.StatusServiceUnavailable, status)
		}
	})
	r.GET("/ready", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ready"})
	})

	// Metrics endpoint (Prometheus)
	if os.Getenv("ENABLE_METRICS") != "false" {
		r.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}

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

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	// Close gRPC connections
	if err := authClient.Close(); err != nil {
		log.Error("Failed to close auth client", "error", err)
	}
	if err := userClient.Close(); err != nil {
		log.Error("Failed to close user client", "error", err)
	}
	if err := tenantClient.Close(); err != nil {
		log.Error("Failed to close tenant client", "error", err)
	}

	log.Info("API Gateway stopped")
}

func getServiceURL(envVar, defaultValue string) string {
	url := os.Getenv(envVar)
	if url == "" {
		return defaultValue
	}
	return url
}
