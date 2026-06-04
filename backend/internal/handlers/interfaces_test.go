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

func TestInterfaceEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create an asset
	asset := models.Asset{
		Name: "Router-01",
		Type: "Router",
	}
	database.DB.Create(&asset)

	var interfaceID string

	t.Run("Create Network Interface", func(t *testing.T) {
		reqBody := handlers.NetworkInterfaceCreateReq{
			Name:       "GigabitEthernet0/1",
			Type:       "ethernet",
			MACAddress: "00:11:22:33:44:55",
			IPAddress:  "10.0.1.1",
			Speed:      "1 Gbps",
			MTU:        1500,
			VLAN:       "VLAN 10",
			Status:     "up",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/interfaces")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset.ID)

		err := handlers.CreateNetworkInterface(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var netInterface models.NetworkInterface
		json.Unmarshal(rec.Body.Bytes(), &netInterface)
		if netInterface.Name != "GigabitEthernet0/1" {
			t.Errorf("Expected name 'GigabitEthernet0/1', got '%s'", netInterface.Name)
		}
		if netInterface.MTU != 1500 {
			t.Errorf("Expected MTU 1500, got %d", netInterface.MTU)
		}
		interfaceID = netInterface.ID
	})

	t.Run("Read Network Interfaces", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/interfaces")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset.ID)

		err := handlers.ReadNetworkInterfaces(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var interfaces []models.NetworkInterface
		json.Unmarshal(rec.Body.Bytes(), &interfaces)
		if len(interfaces) != 1 {
			t.Errorf("Expected 1 interface, got %d", len(interfaces))
		}
	})

	t.Run("Update Network Interface", func(t *testing.T) {
		reqBody := map[string]any{
			"status": "down",
			"mtu":    9000, // Jumbo frames!
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/interfaces/:interface_id")
		c.SetParamNames("interface_id")
		c.SetParamValues(interfaceID)

		err := handlers.UpdateNetworkInterface(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var netInterface models.NetworkInterface
		json.Unmarshal(rec.Body.Bytes(), &netInterface)
		if netInterface.Status != "down" {
			t.Errorf("Expected status 'down', got '%s'", netInterface.Status)
		}
		if netInterface.MTU != 9000 {
			t.Errorf("Expected MTU 9000, got %d", netInterface.MTU)
		}
	})

	t.Run("Delete Network Interface", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/interfaces/:interface_id")
		c.SetParamNames("interface_id")
		c.SetParamValues(interfaceID)

		err := handlers.DeleteNetworkInterface(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})
}
