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

func TestRelationshipEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create some assets
	asset1 := models.Asset{Name: "Switch-A", Type: "Switch"}
	asset2 := models.Asset{Name: "Router-B", Type: "Router"}
	asset3 := models.Asset{Name: "Switch-C", Type: "Switch"}
	database.DB.Create(&asset1)
	database.DB.Create(&asset2)
	database.DB.Create(&asset3)

	t.Run("Sync Relationships", func(t *testing.T) {
		reqBody := handlers.RelationshipSyncReq{
			UplinkIDs:   []string{asset2.ID},
			DownlinkIDs: []string{asset3.ID},
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/relationships")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset1.ID)

		err := handlers.SyncRelationships(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})

	t.Run("Read Relationships", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/assets/:asset_id/relationships")
		c.SetParamNames("asset_id")
		c.SetParamValues(asset1.ID)

		err := handlers.ReadRelationships(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var res map[string][]models.Asset
		json.Unmarshal(rec.Body.Bytes(), &res)
		if len(res["uplinks"]) != 1 {
			t.Errorf("Expected 1 uplink, got %d", len(res["uplinks"]))
		}
		if len(res["downlinks"]) != 1 {
			t.Errorf("Expected 1 downlink, got %d", len(res["downlinks"]))
		}
		if res["uplinks"][0].Name != "Router-B" {
			t.Errorf("Expected uplink to be 'Router-B', got '%s'", res["uplinks"][0].Name)
		}
	})
}
