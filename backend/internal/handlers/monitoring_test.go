package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/handlers"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

func TestMonitoringEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Seed at least one asset
	asset := models.Asset{
		Name:   "Test Server",
		Type:   "Server",
		Status: "active",
	}
	database.DB.Create(&asset)

	t.Run("Get Latency Ping Matrix", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/monitoring/ping", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.GetLatencyPing(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var pings []handlers.LatencyPingRes
		json.Unmarshal(rec.Body.Bytes(), &pings)
		if len(pings) < 1 {
			t.Errorf("Expected at least 1 latency ping result, got %d", len(pings))
		}

		if pings[0].Status != "online" {
			t.Errorf("Expected status 'online', got '%s'", pings[0].Status)
		}
	})
}
