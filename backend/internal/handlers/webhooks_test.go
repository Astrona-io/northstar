package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"cmdb-backend/internal/handlers"
	"cmdb-backend/internal/models"
)

func TestWebhookEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	var webhookID string

	t.Run("Create Webhook", func(t *testing.T) {
		reqBody := handlers.WebhookCreateReq{
			URL:   "https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX",
			Event: "asset:create",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/webhooks", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateWebhook(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var hook models.Webhook
		json.Unmarshal(rec.Body.Bytes(), &hook)
		if hook.Event != "asset:create" {
			t.Errorf("Expected event 'asset:create', got '%s'", hook.Event)
		}
		webhookID = hook.ID
	})

	t.Run("Read Webhooks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/webhooks", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadWebhooks(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})

	t.Run("Delete Webhook", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/webhooks/:id")
		c.SetParamNames("id")
		c.SetParamValues(webhookID)

		err := handlers.DeleteWebhook(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})
}
