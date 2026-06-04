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

func TestUserEndpoints(t *testing.T) {
	setupTestDB()
	e := echo.New()

	// Pre-create user and permission
	perm := models.Permission{Name: "catalog:write", Effect: "deny"}
	database.DB.Create(&perm)

	user := models.User{
		Username: "hybrid-test-spec",
		Password: "password123",
		Role:     "operator",
	}
	database.DB.Create(&user)

	t.Run("Read Users", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := handlers.ReadUsers(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}
	})

	t.Run("Update User Overrides", func(t *testing.T) {
		reqBody := map[string]any{
			"permission_ids": []string{perm.ID},
		}
		b, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader(b))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/api/users/:id")
		c.SetParamNames("id")
		c.SetParamValues(user.ID)

		err := handlers.UpdateUser(c)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if rec.Code != http.StatusOK {
			t.Errorf("Expected status OK, got %d", rec.Code)
		}

		var updated models.User
		json.Unmarshal(rec.Body.Bytes(), &updated)
		if len(updated.Permissions) != 1 {
			t.Errorf("Expected 1 permission override, got %d", len(updated.Permissions))
		}
		if updated.Permissions[0].Name != "catalog:write" {
			t.Errorf("Expected override 'catalog:write', got '%s'", updated.Permissions[0].Name)
		}
	})
}
