package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/handlers"
	"cmdb-backend/internal/models"
)

func TestMaintenanceEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create an asset so we can link maintenance windows to it
	asset := models.Asset{
		Name: "Server to Maintain",
		Type: "Server",
	}
	database.DB.Create(&asset)

	var windowID string

	t.Run("Create Maintenance Window", func(t *testing.T) {
		reqBody := handlers.MaintenanceWindowCreateReq{
			Title:     "Quarterly Maintenance",
			StartTime: time.Now().Add(1 * time.Hour),
			EndTime:   time.Now().Add(2 * time.Hour),
			Status:    "scheduled",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/maintenance")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset.ID)

		err := handlers.CreateMaintenanceWindow(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var window models.MaintenanceWindow
		json.Unmarshal(rec.Body.Bytes(), &window)
		if window.Title != "Quarterly Maintenance" {
			t.Errorf("Expected title 'Quarterly Maintenance', got '%s'", window.Title)
		}
		windowID = window.ID
	})

	t.Run("Read Maintenance Windows", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/maintenance")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset.ID)

		err := handlers.ReadMaintenanceWindows(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var windows []models.MaintenanceWindow
		json.Unmarshal(rec.Body.Bytes(), &windows)
		if len(windows) != 1 {
			t.Errorf("Expected 1 window, got %d", len(windows))
		}
	})

	t.Run("Update Maintenance Window", func(t *testing.T) {
		reqBody := map[string]any{
			"status": "completed",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/maintenance/:window_id")
		c.SetParamNames("window_id")
		c.SetParamValues(windowID)

		err := handlers.UpdateMaintenanceWindow(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var window models.MaintenanceWindow
		json.Unmarshal(rec.Body.Bytes(), &window)
		if window.Status != "completed" {
			t.Errorf("Expected status 'completed', got '%s'", window.Status)
		}
	})

	t.Run("Delete Maintenance Window", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/maintenance/:window_id")
		c.SetParamNames("window_id")
		c.SetParamValues(windowID)

		err := handlers.DeleteMaintenanceWindow(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})
}
