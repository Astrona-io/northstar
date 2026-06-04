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

func TestDiscoveryEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	t.Run("Scan Subnet Sweep", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/discovery/scan?subnet=10.0.1.0/24", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ScanSubnet(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var devices []handlers.DiscoveredDevice
		json.Unmarshal(rec.Body.Bytes(), &devices)
		if len(devices) < 1 {
			t.Errorf("Expected at least 1 discovered active port, got %d", len(devices))
		}
	})

	t.Run("Import Discovered Device", func(t *testing.T) {
		reqBody := handlers.ImportDeviceReq{
			Name:         "Imported-Switch",
			Type:         "Switch",
			IPAddress:    "10.0.1.50",
			Description:  "Cisco switch imported via test sweep",
			Manufacturer: "Cisco Systems",
			Model:        "Catalyst 9300",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/discovery/import", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ImportDevice(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var asset models.Asset
		json.Unmarshal(rec.Body.Bytes(), &asset)
		if asset.Name != "Imported-Switch" {
			t.Errorf("Expected imported asset name to be 'Imported-Switch', got '%s'", asset.Name)
		}
	})
}
