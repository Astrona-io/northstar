package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/handlers"
	"cmdb-backend/internal/models"
)

func setupTestDB() {
	handlers.IsTestMode = true
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(
		&models.Asset{},
		&models.MaintenanceWindow{},
		&models.NetworkInterface{},
		&models.AuditLog{},
		&models.AssetRelationship{},
		&models.Manufacturer{},
		&models.Category{},
		&models.SubGroup{},
		&models.DeviceModel{},
		&models.DeviceModelRevision{},
		&models.Permission{},
		&models.Group{},
		&models.User{},
		&models.DatacenterType{},
		&models.DatacenterFloor{},
		&models.DatacenterWall{},
		&models.Datacenter{},
		&models.Rack{},
		&models.CustomField{},
		&models.Webhook{},
		&models.LicenseAgreement{},
	)
	database.DB = db
}

func TestAssetEndpoints(t *testing.T) {
	var assetID string
	setupTestDB()
	e := echo.New()

	t.Run("Create Asset", func(t *testing.T) {
		reqBody := handlers.AssetCreateReq{
			Name:   "Test Server",
			Type:   "Server",
			Status: "active",
			Properties: models.JSONMap{
				"cpu_model": "Intel Xeon",
				"ignored":   "should be removed",
			},
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/assets/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateAsset(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var asset models.Asset
		json.Unmarshal(rec.Body.Bytes(), &asset)
		assetID = asset.ID
		t.Logf("Response body: %s", rec.Body.String())
		t.Logf("Asset Properties: %v", asset.Properties)
		if asset.Name != "Test Server" {
			t.Errorf("Expected name 'Test Server', got '%s'", asset.Name)
		}
		if asset.Properties["cpu_model"] != "Intel Xeon" {
			t.Errorf("Expected cpu_model 'Intel Xeon', got '%v'", asset.Properties["cpu_model"])
		}
		if _, ok := asset.Properties["ignored"]; ok {
			t.Errorf("Expected ignored property to be removed")
		}
	})

	t.Run("Read Assets", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/assets/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadAssets(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var assets []models.Asset
		json.Unmarshal(rec.Body.Bytes(), &assets)
		if len(assets) != 1 {
			t.Errorf("Expected 1 asset, got %d", len(assets))
		}
	})

	t.Run("Read Assets with Search and Filter", func(t *testing.T) {
		// First verify matching search works
		req := httptest.NewRequest(http.MethodGet, "/api/assets/?search=Server&type=Server&status=active", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadAssets(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var assets []models.Asset
		json.Unmarshal(rec.Body.Bytes(), &assets)
		if len(assets) != 1 {
			t.Errorf("Expected 1 asset matching criteria, got %d", len(assets))
		}

		// Now verify non-matching criteria returns 0 results
		req2 := httptest.NewRequest(http.MethodGet, "/api/assets/?search=Database&type=Database&status=active", nil)
		rec2 := httptest.NewRecorder()
		c2 := e.NewContext(req2, rec2)

		err2 := handlers.ReadAssets(c2)
		if err2 != nil {
			t.Fatalf("Expected no error, got %v", err2)
		}
		var emptyAssets []models.Asset
		json.Unmarshal(rec2.Body.Bytes(), &emptyAssets)
		if len(emptyAssets) != 0 {
			t.Errorf("Expected 0 assets for non-matching search, got %d", len(emptyAssets))
		}
	})

	t.Run("Update Asset", func(t *testing.T) {
		reqBody := map[string]any{
			"status": "maintenance",
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id")
		c.SetParamNames("asset_id")
		c.SetParamValues(assetID)

		err := handlers.UpdateAsset(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var asset models.Asset
		json.Unmarshal(rec.Body.Bytes(), &asset)
		assetID = asset.ID
		if asset.Status != "maintenance" {
			t.Errorf("Expected status 'maintenance', got '%s'", asset.Status)
		}
	})

	t.Run("Delete Asset", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id")
		c.SetParamNames("asset_id")
		c.SetParamValues(assetID)

		err := handlers.DeleteAsset(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})
}
