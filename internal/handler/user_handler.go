package handler

import (
"go.uber.org/zap"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhvcorp/go-shared/logger"
	"github.com/vhvcorp/go-api-gateway/internal/client"
)

// UserHandler handles user-related requests
type UserHandler struct {
	client *client.UserClient
	log    *logger.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(client *client.UserClient, log *logger.Logger) *UserHandler {
	return &UserHandler{
		client: client,
		log:    log,
	}
}

// CreateUser forwards create user requests
func (h *UserHandler) CreateUser(c *gin.Context) {
	h.forwardRequest(c, "http://user-service:8082/api/v1/users", "POST")
}

// GetUsers forwards get users requests
func (h *UserHandler) GetUsers(c *gin.Context) {
	url := fmt.Sprintf("http://user-service:8082/api/v1/users?%s", c.Request.URL.RawQuery)
	h.forwardRequest(c, url, "GET")
}

// GetUser forwards get user by ID requests
func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://user-service:8082/api/v1/users/%s", id)
	h.forwardRequest(c, url, "GET")
}

// UpdateUser forwards update user requests
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://user-service:8082/api/v1/users/%s", id)
	h.forwardRequest(c, url, "PUT")
}

// DeleteUser forwards delete user requests
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("http://user-service:8082/api/v1/users/%s", id)
	h.forwardRequest(c, url, "DELETE")
}

// SearchUsers forwards search user requests
func (h *UserHandler) SearchUsers(c *gin.Context) {
	url := fmt.Sprintf("http://user-service:8082/api/v1/users/search?%s", c.Request.URL.RawQuery)
	h.forwardRequest(c, url, "GET")
}

// forwardRequest is a helper method to forward requests
func (h *UserHandler) forwardRequest(c *gin.Context, targetURL, method string) {
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		h.log.Error("Failed to create request", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	req.Header.Set("X-Correlation-ID", c.GetString("correlation_id"))
	req.Header.Set("X-Tenant-ID", c.GetString("tenant_id"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Error("Failed to forward request", zap.Error(err), zap.String("url", targetURL))
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Error("Failed to read response", zap.Error(err))
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
