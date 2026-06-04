package handlers

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type WebhookCreateReq struct {
	URL   string `json:"url"`
	Event string `json:"event"` // e.g. "asset:create", "asset:update", "asset:delete"
}

// CreateWebhook handles POST /api/webhooks
func CreateWebhook(c echo.Context) error {
	var req WebhookCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if req.URL == "" || req.Event == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "URL and Event fields are required"})
	}

	webhook := models.Webhook{
		URL:   req.URL,
		Event: req.Event,
	}

	if err := database.DB.Create(&webhook).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to register webhook: " + err.Error()})
	}

	return c.JSON(http.StatusOK, webhook)
}

// ReadWebhooks handles GET /api/webhooks
func ReadWebhooks(c echo.Context) error {
	var list []models.Webhook
	database.DB.Order("id desc").Find(&list)
	return c.JSON(http.StatusOK, list)
}

// DeleteWebhook handles DELETE /api/webhooks/{id}
func DeleteWebhook(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid webhook ID"})
	}

	var webhook models.Webhook
	if err := database.DB.First(&webhook, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Webhook not found"})
	}

	if err := database.DB.Delete(&webhook).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete webhook: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Webhook deleted successfully"})
}

// DispatchWebhookEvent broadcasts an event payload asynchronously to all registered webhook endpoints. (Phase 3)
func DispatchWebhookEvent(event string, payload any) {
	var list []models.Webhook
	// Query matching webhooks in DB
	database.DB.Where("event = ?", event).Find(&list)

	if len(list) == 0 {
		return
	}

	log.Printf("[Webhooks Engine] Dispatching event '%s' to %d targets...", event, len(list))

	// Dispatch each callback asynchronously
	for _, hook := range list {
		go func(url string, ev string, data any) {
			b, err := json.Marshal(echo.Map{
				"event":     ev,
				"timestamp": time.Now().Format(time.RFC3339),
				"data":      data,
			})
			if err != nil {
				return
			}

			client := &http.Client{Timeout: 5 * time.Second}
			resp, err := client.Post(url, "application/json", bytes.NewReader(b))
			if err != nil {
				log.Printf("[Webhooks Error] Failed to post callback to %s: %v", url, err)
				return
			}
			resp.Body.Close()
			log.Printf("[Webhooks Engine] Successfully delivered event '%s' to %s (Status: %s)", ev, url, resp.Status)
		}(hook.URL, hook.Event, payload)
	}
}
