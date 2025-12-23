package handler

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"github.com/longvhv/saas-framework-go/services/api-gateway/internal/client"
)

// AuthHandler handles auth-related requests
type AuthHandler struct {
	client *client.AuthClient
	log    *logger.Logger
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(client *client.AuthClient, log *logger.Logger) *AuthHandler {
	return &AuthHandler{
		client: client,
		log:    log,
	}
}

// Register forwards register requests to auth service
func (h *AuthHandler) Register(c *gin.Context) {
	h.forwardRequest(c, "http://auth-service:8081/api/v1/auth/register", "POST")
}

// Login forwards login requests to auth service
func (h *AuthHandler) Login(c *gin.Context) {
	h.forwardRequest(c, "http://auth-service:8081/api/v1/auth/login", "POST")
}

// Logout forwards logout requests to auth service
func (h *AuthHandler) Logout(c *gin.Context) {
	h.forwardRequest(c, "http://auth-service:8081/api/v1/auth/logout", "POST")
}

// RefreshToken forwards refresh token requests to auth service
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	h.forwardRequest(c, "http://auth-service:8081/api/v1/auth/refresh", "POST")
}

// forwardRequest is a helper method to forward requests to backend services
func (h *AuthHandler) forwardRequest(c *gin.Context, targetURL, method string) {
	// Read request body
	var bodyBytes []byte
	if c.Request.Body != nil {
		bodyBytes, _ = io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	// Create new request
	req, err := http.NewRequest(method, targetURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		h.log.Error("Failed to create request", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to forward request"})
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", c.GetHeader("Content-Type"))
	req.Header.Set("X-Correlation-ID", c.GetString("correlation_id"))

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.log.Error("Failed to forward request", "error", err, "url", targetURL)
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		h.log.Error("Failed to read response", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	// Parse and return response
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		h.log.Error("Failed to parse response", "error", err)
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), respBody)
		return
	}

	c.JSON(resp.StatusCode, result)
}
