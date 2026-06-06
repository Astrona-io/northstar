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

	t.Run("Create, Update, and Pin Device Model Revision", func(t *testing.T) {
		// 1. Create a Device Model with Ports
		reqBody := handlers.DeviceModelCreateReq{
			ManufacturerID: m.ID,
			ModelName:      "Router-X",
			CategoryIDs:    []string{},
			GeneralInfo:    "High-speed core router",
			Ports:          models.JSONMap{"RJ45": float64(8), "SFP+": float64(2)},
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
			t.Fatalf("Expected status OK, got %d", rec.Code)
		}

		var device models.DeviceModel
		json.Unmarshal(rec.Body.Bytes(), &device)
		if device.Revision != 1 {
			t.Errorf("Expected initial revision 1, got %d", device.Revision)
		}
		if val, ok := device.Ports["SFP+"].(float64); !ok || val != 2 {
			t.Errorf("Expected 2 SFP+ ports, got %v", device.Ports["SFP+"])
		}

		// Verify first revision record exists
		var revRecord models.DeviceModelRevision
		if err := database.DB.First(&revRecord, "device_model_id = ? AND revision = 1", device.ID).Error; err != nil {
			t.Fatalf("Expected revision 1 record to exist in DB, got error: %v", err)
		}

		// 2. Create an Asset referencing this Device Model (should auto-pin to Revision 1)
		asset := models.Asset{
			Name:          "Router-01",
			Type:          "Network",
			Status:        "active",
			DeviceModelID: &device.ID,
		}
		if err := database.DB.Create(&asset).Error; err != nil {
			t.Fatalf("Expected to create asset successfully, got %v", err)
		}

		// Simulate CreateAsset resolution
		assetCreateReq := handlers.AssetCreateReq{
			Name:          "Router-01-API",
			Type:          "Network",
			Status:        "active",
			DeviceModelID: &device.ID,
		}
		bAsset, _ := json.Marshal(assetCreateReq)
		reqAsset := httptest.NewRequest(http.MethodPost, "/api/assets/", bytes.NewReader(bAsset))
		reqAsset.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recAsset := httptest.NewRecorder()
		cAsset := e.NewContext(reqAsset, recAsset)
		if err := handlers.CreateAsset(cAsset); err != nil {
			t.Fatalf("Expected no error creating asset via API, got %v", err)
		}
		var createdAsset models.Asset
		json.Unmarshal(recAsset.Body.Bytes(), &createdAsset)
		if createdAsset.DeviceModelRevision == nil || *createdAsset.DeviceModelRevision != 1 {
			t.Errorf("Expected asset to pin to DeviceModelRevision 1, got %v", createdAsset.DeviceModelRevision)
		}

		// 3. Update the Device Model (since it is used in assets, it should bump revision to 2)
		updateBody := map[string]any{
			"model_name": "Router-X-Updated",
			"ports":      models.JSONMap{"RJ45": float64(8), "SFP+": float64(2), "QSFP+": float64(1)},
		}
		bUp, _ := json.Marshal(updateBody)
		reqUp := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(bUp))
		reqUp.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		recUp := httptest.NewRecorder()
		cUp := e.NewContext(reqUp, recUp)
		cUp.SetPath("/api/devices/:device_id")
		cUp.SetParamNames("device_id")
		cUp.SetParamValues(device.ID)

		if err := handlers.UpdateDeviceModel(cUp); err != nil {
			t.Fatalf("Expected no error updating device model, got %v", err)
		}
		if recUp.Code != http.StatusOK {
			t.Fatalf("Expected update status OK, got %d, body: %s", recUp.Code, recUp.Body.String())
		}

		var updatedDevice models.DeviceModel
		json.Unmarshal(recUp.Body.Bytes(), &updatedDevice)
		if updatedDevice.Revision != 2 {
			t.Errorf("Expected revision bumped to 2, got %d", updatedDevice.Revision)
		}

		// Verify revision 2 record exists
		var revRecord2 models.DeviceModelRevision
		if err := database.DB.First(&revRecord2, "device_model_id = ? AND revision = 2", device.ID).Error; err != nil {
			t.Fatalf("Expected revision 2 record to exist in DB, got error: %v", err)
		}

		// 4. Verify existing asset still points to Revision 1 (not automatically bumped!)
		var reloadedAsset models.Asset
		if err := database.DB.First(&reloadedAsset, "id = ?", createdAsset.ID).Error; err != nil {
			t.Fatalf("Expected to reload asset, got %v", err)
		}
		if reloadedAsset.DeviceModelRevision == nil || *reloadedAsset.DeviceModelRevision != 1 {
			t.Errorf("Expected existing asset to remain pinned to Revision 1, got %v", reloadedAsset.DeviceModelRevision)
		}

		// Read asset details via ReadAsset API to confirm preloaded revision details are Revision 1
		reqRead := httptest.NewRequest(http.MethodGet, "/", nil)
		recRead := httptest.NewRecorder()
		cRead := e.NewContext(reqRead, recRead)
		cRead.SetPath("/api/assets/:asset_id")
		cRead.SetParamNames("asset_id")
		cRead.SetParamValues(createdAsset.ID)
		if err := handlers.ReadAsset(cRead); err != nil {
			t.Fatalf("Expected no error on ReadAsset, got %v", err)
		}
		var readAsset models.Asset
		json.Unmarshal(recRead.Body.Bytes(), &readAsset)
		if readAsset.DeviceModelRevisionDetails == nil {
			t.Fatalf("Expected DeviceModelRevisionDetails to be populated in ReadAsset")
		}
		if readAsset.DeviceModelRevisionDetails.Revision != 1 {
			t.Errorf("Expected preloaded revision details to be for Revision 1, got %d", readAsset.DeviceModelRevisionDetails.Revision)
		}
	})
}
