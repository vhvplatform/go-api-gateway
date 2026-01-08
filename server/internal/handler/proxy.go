package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
)

// ProxyHandler handles reverse proxying to other services
type ProxyHandler struct {
	log *logger.Logger
}

func NewProxyHandler(log *logger.Logger) *ProxyHandler {
	return &ProxyHandler{log: log}
}

// APIProxy forwards requests to Go microservices
func (h *ProxyHandler) APIProxy(c *gin.Context) {
	// Path format: /api/service-name/api-path
	serviceName := c.Param("service")
	targetHost := h.getServiceURL(serviceName)

	if targetHost == "" {
		h.log.Warn("Unknown API service", zap.String("service", serviceName))
		h.handleFailover(c, "Service not found")
		return
	}

	h.proxyRequest(c, targetHost)
}

// PageProxy forwards requests to React Frontends
func (h *ProxyHandler) PageProxy(c *gin.Context) {
	// Path format: /page/service-name/page-path
	parts := strings.Split(strings.TrimPrefix(c.Request.URL.Path, "/"), "/")
	if len(parts) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page path"})
		return
	}

	serviceName := parts[1]
	targetHost := h.getServiceURL(serviceName + "-ui") // Convention: service-name-ui

	if targetHost == "" {
		h.handleFailover(c, "Page service not found")
		return
	}

	h.proxyRequest(c, targetHost)
}

// UploadProxy forwards requests to file-service
func (h *ProxyHandler) UploadProxy(c *gin.Context) {
	// Path format: /upload/file-key
	targetHost := h.getServiceURL("file-service")

	if targetHost == "" {
		h.handleFailover(c, "Upload service not found")
		return
	}

	h.proxyRequest(c, targetHost)
}

// SlugProxy handles pretty URLs (slugs)
func (h *ProxyHandler) SlugProxy(c *gin.Context) {
	// Fallback for any other path
	// Redirect to default service of tenant if available
	tenantDefault := c.GetString("tenant_default_service")
	if tenantDefault == "" {
		tenantDefault = "cms-service" // System fallback
	}

	targetHost := h.getServiceURL(tenantDefault)
	if targetHost == "" {
		h.handleFailover(c, "Slug handler not found")
		return
	}
	h.proxyRequest(c, targetHost)
}

func (h *ProxyHandler) proxyRequest(c *gin.Context, target string) {
	targetURL, err := url.Parse(target)
	if err != nil {
		h.log.Error("Failed to parse target URL", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// Director to modify request
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Host = targetURL.Host
		req.URL.Scheme = targetURL.Scheme
		req.URL.Host = targetURL.Host
	}

	// ErrorHandler for Failover
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		h.log.Error("Proxy error, triggering failover", zap.Error(err), zap.String("target", target))
		h.handleFailover(c, err.Error())
	}

	// ServeHTTP
	proxy.ServeHTTP(c.Writer, c.Request)
}

func (h *ProxyHandler) handleFailover(c *gin.Context, originalErr string) {
	// 1. Try Tenant Default Service (from context or db)
	tenantDefault := c.GetString("tenant_default_service")
	if tenantDefault != "" {
		h.log.Info("Failing over to tenant default", zap.String("service", tenantDefault))
		c.Set("tenant_default_service", "") // Clear to avoid infinite loop
		h.proxyRequest(c, h.getServiceURL(tenantDefault))
		return
	}

	// 2. Try System Default Service (e.g. Dashboard or Login)
	systemDefault := "dashboard-service"
	target := h.getServiceURL(systemDefault)
	if target != "" {
		h.log.Info("Failing over to system default", zap.String("service", systemDefault))
		h.proxyRequest(c, target)
		return
	}

	// 3. Last resort: Auth login
	h.log.Info("Redirecting to auth login")
	c.Redirect(http.StatusFound, "/auth/login?error="+url.QueryEscape(originalErr))
}

func (h *ProxyHandler) getServiceURL(serviceName string) string {
	// In a real scenario, this would use Service Discovery (consul, k8s dns, etc.)
	// For now, mapping from env or simple convention
	envVar := strings.ToUpper(strings.ReplaceAll(serviceName, "-", "_")) + "_URL"
	url := os.Getenv(envVar)
	if url != "" {
		return url
	}

	// Fallback to convention: http://service-name:port
	// This is a simplified example
	return ""
}
