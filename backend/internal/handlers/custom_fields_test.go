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

func TestCustomFieldsEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	var fieldID string

	t.Run("Create Custom Field Specification", func(t *testing.T) {
		reqBody := handlers.CustomFieldCreateReq{
			Name:       "warranty_expiration_date",
			FieldType:  "text",
			AssetType:  "Server",
			IsRequired: true,
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/custom-fields", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.CreateCustomField(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var field models.CustomField
		json.Unmarshal(rec.Body.Bytes(), &field)
		if field.Name != "warranty_expiration_date" {
			t.Errorf("Expected custom field name 'warranty_expiration_date', got '%s'", field.Name)
		}
		fieldID = field.ID
	})

	t.Run("Read Custom Fields by Asset Type", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/custom-fields?asset_type=Server", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadCustomFields(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var fields []models.CustomField
		json.Unmarshal(rec.Body.Bytes(), &fields)
		if len(fields) != 2 {
			t.Errorf("Expected 2 custom fields matching criteria, got %d", len(fields))
		}
	})

	t.Run("Delete Custom Field", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/custom-fields/:id")
		c.SetParamNames("id")
		c.SetParamValues(fieldID)

		err := handlers.DeleteCustomField(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})
}
