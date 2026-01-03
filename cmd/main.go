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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/vhvplatform/go-api-gateway/internal/cache"
	"github.com/vhvplatform/go-api-gateway/internal/circuitbreaker"
	"github.com/vhvplatform/go-api-gateway/internal/client"
	"github.com/vhvplatform/go-api-gateway/internal/handler"
	"github.com/vhvplatform/go-api-gateway/internal/health"
	"github.com/vhvplatform/go-api-gateway/internal/router"
	"github.com/vhvplatform/go-api-gateway/internal/tracing"
	"github.com/vhvplatform/go-shared/config"
	"github.com/vhvplatform/go-shared/logger"
	pkgmiddleware "github.com/vhvplatform/go-shared/middleware"
	"go.uber.org/zap"

	_ "github.com/vhvplatform/go-api-gateway/docs"
)

// @title API Gateway
// @version 1.0
// @description Unified API Gateway for VHV Platform Microservices
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.vhvplatform.com/support
// @contact.email support@vhvplatform.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Initialize logger
	log := logger.NewLogger()
	defer log.Sync()

	log.Info("Starting API Gateway...")

	// Create main context for graceful shutdown
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize distributed tracing (optional)
	if os.Getenv("ENABLE_TRACING") == "true" {
		jaegerURL := getServiceURL("JAEGER_URL", "http://jaeger:14268/api/traces")
		tp, err := tracing.InitTracer("api-gateway", jaegerURL)
		if err != nil {
			log.Error("Failed to initialize tracer", zap.Error(err))
		} else {
			defer func() {
				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := tp.Shutdown(ctx); err != nil {
					log.Error("Failed to shutdown tracer", zap.Error(err))
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
			log.Error("Failed to initialize cache", zap.Error(err))
		} else {
			defer cacheClient.Close()
			log.Info("Redis cache initialized")
		}
	}

	// Initialize circuit breaker
	_ = circuitbreaker.NewCircuitBreaker()
	// TODO: Use circuit breaker in handlers when making external calls

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
	r.Use(pkgmiddleware.Recovery(log))

	// Correlation ID middleware (should be first)
	r.Use(pkgmiddleware.CorrelationID())

	// Logging middleware
	r.Use(pkgmiddleware.Logger(log))

	// Metrics middleware (if enabled)
	if os.Getenv("ENABLE_METRICS") != "false" {
		r.Use(pkgmiddleware.DefaultMetrics("api_gateway"))
	}

	// Compression middleware
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	// Request validation middleware
	r.Use(pkgmiddleware.RequestValidation())

	// Request size limit middleware
	maxRequestSize := int64(10485760) // 10MB default
	if size := os.Getenv("MAX_REQUEST_SIZE"); size != "" {
		if parsedSize, err := strconv.ParseInt(size, 10, 64); err == nil {
			maxRequestSize = parsedSize
		}
	}
	r.Use(pkgmiddleware.RequestSizeLimit(maxRequestSize))

	// Timeout middleware
	r.Use(pkgmiddleware.Timeout(30 * time.Second))

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
	r.Use(pkgmiddleware.PerIP(rateLimit, rateBurst))

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

	// Swagger endpoint
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
		log.Info("API Gateway started", zap.String("port", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API Gateway...")

	// Cancel context to stop background goroutines
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Close gRPC connections
	if err := authClient.Close(); err != nil {
		log.Error("Failed to close auth client", zap.Error(err))
	}
	if err := userClient.Close(); err != nil {
		log.Error("Failed to close user client", zap.Error(err))
	}
	if err := tenantClient.Close(); err != nil {
		log.Error("Failed to close tenant client", zap.Error(err))
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
