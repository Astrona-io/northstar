package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/handlers"
	"cmdb-backend/internal/models"
)

func TestAuditLogEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create an asset which will automatically trigger the "create" audit hook/log
	asset := models.Asset{
		Name: "Database Server",
		Type: "Database",
	}
	database.DB.Create(&asset)

	// Since we used GORM directly to pre-create, let's manually record a log or trigger a handler
	handlers.RecordAuditLog(asset.ID, "create", "Manual trigger check")

	t.Run("Read Audit Logs", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/logs")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset.ID)

		err := handlers.ReadAuditLogs(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var logs []models.AuditLog
		json.Unmarshal(rec.Body.Bytes(), &logs)
		if len(logs) != 1 {
			t.Errorf("Expected 1 audit log trace, got %d", len(logs))
		}
		if logs[0].Action != "create" {
			t.Errorf("Expected action 'create', got '%s'", logs[0].Action)
		}
	})
}
