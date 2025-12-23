package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/longvhv/saas-framework-go/pkg/httpclient"
	"github.com/longvhv/saas-framework-go/pkg/logger"
	"go.uber.org/zap"
)

// NotificationHandler handles notification-related requests
type NotificationHandler struct {
	baseURL string
	client  *httpclient.Client
	log     *logger.Logger
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(baseURL string, log *logger.Logger) *NotificationHandler {
	client := httpclient.NewClient(
		httpclient.WithBaseURL(baseURL),
		httpclient.WithRetry(3, 1),
		httpclient.WithCircuitBreaker(),
	)
	
	return &NotificationHandler{
		baseURL: baseURL,
		client:  client,
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
	var body map[string]interface{}
	if c.Request.Body != nil && method != "GET" {
		if err := c.ShouldBindJSON(&body); err != nil {
			h.log.Error("Failed to parse request body", zap.Error(err))
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}
	}

	// Extract path from full URL
	path := targetURL
	if len(h.baseURL) > 0 && len(targetURL) > len(h.baseURL) {
		path = targetURL[len(h.baseURL):]
	}

	var result map[string]interface{}
	var err error

	switch method {
	case "GET":
		err = h.client.Get(c.Request.Context(), path, &result)
	case "POST":
		err = h.client.Post(c.Request.Context(), path, body, &result)
	case "PUT":
		err = h.client.Put(c.Request.Context(), path, body, &result)
	case "DELETE":
		err = h.client.Delete(c.Request.Context(), path)
		if err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "deleted"})
			return
		}
	default:
		h.log.Error("Unsupported HTTP method", zap.String("method", method))
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "Method not allowed"})
		return
	}

	if err != nil {
		h.log.Error("Failed to forward request", zap.Error(err), zap.String("url", targetURL))
		c.JSON(http.StatusBadGateway, gin.H{"error": "Service unavailable"})
		return
	}

	c.JSON(http.StatusOK, result)
}
