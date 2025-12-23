package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/client"
)

// TenantHandler handles tenant-related requests
type TenantHandler struct {
	client *client.TenantClient
	log    *logger.Logger
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(client *client.TenantClient, log *logger.Logger) *TenantHandler {
	return &TenantHandler{
		client: client,
		log:    log,
	}
}

// CreateTenant forwards create tenant requests
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	h.forwardRequest(c, "http://tenant-service:8083/api/v1/tenants", "POST")
}

// GetTenants forwards get tenants requests
func (h *TenantHandler) GetTenants(c *gin.Context) {
	url := fmt.Sprintf("http://tenant-service:8083/api/v1/tenants?%s", c.Request.URL.RawQuery)
	h.forwardRequest(c, url, "GET")
}

// GetTenant forwards get tenant by ID requests
func (h *TenantHandler) GetTenant(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://tenant-service:8083/api/v1/tenants/%s", id)
	h.forwardRequest(c, url, "GET")
}

// UpdateTenant forwards update tenant requests
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://tenant-service:8083/api/v1/tenants/%s", id)
	h.forwardRequest(c, url, "PUT")
}

// DeleteTenant forwards delete tenant requests
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://tenant-service:8083/api/v1/tenants/%s", id)
	h.forwardRequest(c, url, "DELETE")
}

// AddUserToTenant forwards add user to tenant requests
func (h *TenantHandler) AddUserToTenant(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://tenant-service:8083/api/v1/tenants/%s/users", id)
	h.forwardRequest(c, url, "POST")
}

// RemoveUserFromTenant forwards remove user from tenant requests
func (h *TenantHandler) RemoveUserFromTenant(c *gin.Context) {
	id := c.Param("id")
	userID := c.Param("user_id")
	url := fmt.Sprintf("http://tenant-service:8083/api/v1/tenants/%s/users/%s", id, userID)
	h.forwardRequest(c, url, "DELETE")
}

// forwardRequest is a helper method to forward requests
func (h *TenantHandler) forwardRequest(c *gin.Context, targetURL, method string) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		h.log.Error("Failed to create request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	req.Header.Set("X-Correlation-ID", c.GetString("correlation_id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Error("Failed to forward request", "error", err, "url", targetURL)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Error("Failed to read response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
		return
	}

	c.JSON(resp.StatusCode, result)
}
