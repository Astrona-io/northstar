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

func TestDCIMEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create datacenter and rack
	dc := models.Datacenter{Name: "Copenhagen-Mock", Type: "homelab", Country: "Denmark", City: "Copenhagen"}
	database.DB.Create(&dc)

	rack := models.Rack{Name: "Lab-Rack-01", DatacenterID: dc.ID, HeightU: 42}
	database.DB.Create(&rack)

	t.Run("Read Datacenters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/datacenters", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadDatacenters(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var dcs []models.Datacenter
		json.Unmarshal(rec.Body.Bytes(), &dcs)
		if len(dcs) != 1 {
			t.Errorf("Expected 1 datacenter, got %d", len(dcs))
		}
	})

	t.Run("Read Racks", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/racks", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadRacks(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})

	t.Run("Read Rack Details with preloaded assets", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/racks/:id")
		c.SetParamNames("id")
		c.SetParamValues(rack.ID)

		err := handlers.ReadRack(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var r models.Rack
		json.Unmarshal(rec.Body.Bytes(), &r)
		if r.Name != "Lab-Rack-01" {
			t.Errorf("Expected rack name 'Lab-Rack-01', got '%s'", r.Name)
		}
	})

	t.Run("Create Datacenter", func(t *testing.T) {
		bodyBytes := []byte(`{"name":"New-DC-Test","type":"homelab","country":"Denmark","city":"Copenhagen","properties":{"uplink_speed":"1 Gbps","public_ip":"127.0.0.1"}}`)
		req := httptest.NewRequest(http.MethodPost, "/api/datacenters", bytes.NewReader(bodyBytes))
		// Send JSON request body
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateDatacenter(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d, body: %s", rec.Code, rec.Body.String())
		}

		var created models.Datacenter
		json.Unmarshal(rec.Body.Bytes(), &created)
		if created.Name != "New-DC-Test" {
			t.Errorf("Expected name 'New-DC-Test', got '%s'", created.Name)
		}
	})
}
