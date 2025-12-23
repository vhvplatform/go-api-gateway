package handler

import (
"go.uber.org/zap"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/logger"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	baseURL string
	log     *logger.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(baseURL string, log *logger.Logger) *NotificationHandler {
	return &NotificationHandler{
		baseURL: baseURL,
		log:     log,
	}
}

// SendEmail forwards send email requests
func (h *NotificationHandler) SendEmail(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/notifications/email", h.baseURL)
	h.forwardRequest(c, url, "POST")
}

// SendWebhook forwards send webhook requests
func (h *NotificationHandler) SendWebhook(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/notifications/webhook", h.baseURL)
	h.forwardRequest(c, url, "POST")
}

// GetNotifications forwards get notifications requests
func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	url := fmt.Sprintf("%s/api/v1/notifications?%s", h.baseURL, c.Request.URL.RawQuery)
	h.forwardRequest(c, url, "GET")
}

// GetNotification forwards get notification by ID requests
func (h *NotificationHandler) GetNotification(c *gin.Context) {
	id := c.Param("id")
	url := fmt.Sprintf("%s/api/v1/notifications/%s", h.baseURL, id)
	h.forwardRequest(c, url, "GET")
}

// forwardRequest is a helper method to forward requests
func (h *NotificationHandler) forwardRequest(c *gin.Context, targetURL, method string) {
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
