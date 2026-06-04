package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/handlers"
	"cmdb-backend/internal/models"
)

func TestDeviceEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create a Manufacturer brand (strict validation prerequisite - Phase 2 Catalog)
	m := models.Manufacturer{Name: "Cisco Systems", Website: "https://cisco.com"}
	database.DB.Create(&m)

	// Pre-create standard global category fallback
	globalCat := models.Category{Name: "global"}
	database.DB.Create(&globalCat)

	var deviceID string

	t.Run("Create Device Model", func(t *testing.T) {
		reqBody := handlers.DeviceModelCreateReq{
			ManufacturerID: m.ID,
			ModelName:      "Catalyst 9300-48UX",
			CategoryIDs:    []string{}, // defaults to "global"
			GeneralInfo:    "48-Port PoE stackable switch",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/devices/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateDeviceModel(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var device models.DeviceModel
		json.Unmarshal(rec.Body.Bytes(), &device)
		if device.Manufacturer.Name != "Cisco Systems" {
			t.Errorf("Expected preloaded manufacturer name 'Cisco Systems', got '%s'", device.Manufacturer.Name)
		}
		deviceID = device.ID
	})

	t.Run("Read Device Models", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/devices/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadDeviceModels(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var devices []models.DeviceModel
		json.Unmarshal(rec.Body.Bytes(), &devices)
		if len(devices) != 1 {
			t.Errorf("Expected 1 device, got %d", len(devices))
		}
	})

	t.Run("Update Device Model", func(t *testing.T) {
		reqBody := map[string]any{
			"general_info": "Updated info details",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/devices/:device_id")
		c.SetParamNames("device_id")
		c.SetParamValues(deviceID)

		err := handlers.UpdateDeviceModel(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var device models.DeviceModel
		json.Unmarshal(rec.Body.Bytes(), &device)
		if device.GeneralInfo != "Updated info details" {
			t.Errorf("Expected status 'Updated info details', got '%s'", device.GeneralInfo)
		}
	})

	t.Run("Delete Device Model", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/devices/:device_id")
		c.SetParamNames("device_id")
		c.SetParamValues(deviceID)

		err := handlers.DeleteDeviceModel(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})
}
