package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vhvplatform/go-api-gateway/internal/client"
	"github.com/vhvplatform/go-shared/logger"
	"go.uber.org/zap"
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
	var req client.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.client.Register(c.Request.Context(), &req)
	if err != nil {
		h.log.Error("Failed to register", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Registration failed"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Login forwards login requests to auth service
func (h *AuthHandler) Login(c *gin.Context) {
	var req client.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	resp, err := h.client.Login(c.Request.Context(), &req)
	if err != nil {
		h.log.Error("Failed to login", zap.Error(err))
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Logout forwards logout requests to auth service
func (h *AuthHandler) Logout(c *gin.Context) {
	// gRPC Logout usually needs just token or user ID
	// For now mock success
	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

// RefreshToken forwards refresh token requests to auth service
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Implement gRPC call when available
	c.JSON(http.StatusOK, gin.H{"access_token": "new-opaque-token", "expires_in": 3600})
}
