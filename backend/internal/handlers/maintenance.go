package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

// MaintenanceWindowCreateReq defines the payload structure for scheduling maintenance.
type MaintenanceWindowCreateReq struct {
	Title       string    `json:"title"`
	Description *string   `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Status      string    `json:"status"`
}

// CreateMaintenanceWindow handles POST /api/assets/{asset_id}/maintenance
func CreateMaintenanceWindow(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Verify asset exists first
	var asset models.Asset
	if err := database.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	var req MaintenanceWindowCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}

	if req.Title == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Title field is required"})
	}
	if req.Status == "" {
		req.Status = "scheduled"
	}

	window := models.MaintenanceWindow{
		AssetID:     assetID,
		Title:       req.Title,
		Description: req.Description,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Status:      req.Status,
	}

	if err := database.DB.Create(&window).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create maintenance window: " + err.Error()})
	}

	return c.JSON(http.StatusOK, window)
}

// ReadMaintenanceWindows handles GET /api/assets/{asset_id}/maintenance
func ReadMaintenanceWindows(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Verify asset exists first
	var asset models.Asset
	if err := database.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	var windows []models.MaintenanceWindow
	if err := database.DB.Where("asset_id = ?", assetID).Find(&windows).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch maintenance windows: " + err.Error()})
	}

	return c.JSON(http.StatusOK, windows)
}

// UpdateMaintenanceWindow handles PUT /api/maintenance/{window_id}
func UpdateMaintenanceWindow(c echo.Context) error {
	id := c.Param("window_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid maintenance window ID"})
	}

	var window models.MaintenanceWindow
	if err := database.DB.First(&window, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Maintenance window not found"})
	}

	var body map[string]any
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON: " + err.Error()})
	}

	updates := make(map[string]any)
	if val, exists := body["title"]; exists {
		if s, ok := val.(string); ok {
			updates["title"] = s
		}
	}
	if val, exists := body["description"]; exists {
		if val == nil {
			updates["description"] = nil
		} else if s, ok := val.(string); ok {
			updates["description"] = s
		}
	}
	if val, exists := body["start_time"]; exists {
		if s, ok := val.(string); ok {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid start_time format (expected RFC3339): " + err.Error()})
			}
			updates["start_time"] = t
		}
	}
	if val, exists := body["end_time"]; exists {
		if s, ok := val.(string); ok {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid end_time format (expected RFC3339): " + err.Error()})
			}
			updates["end_time"] = t
		}
	}
	if val, exists := body["status"]; exists {
		if s, ok := val.(string); ok {
			updates["status"] = s
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&window).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update maintenance window: " + err.Error()})
		}
	}

	if err := database.DB.First(&window, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to reload updated window"})
	}

	return c.JSON(http.StatusOK, window)
}

// DeleteMaintenanceWindow handles DELETE /api/maintenance/{window_id}
func DeleteMaintenanceWindow(c echo.Context) error {
	id := c.Param("window_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid maintenance window ID"})
	}

	var window models.MaintenanceWindow
	if err := database.DB.First(&window, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Maintenance window not found"})
	}

	if err := database.DB.Delete(&window).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete maintenance window: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Maintenance window deleted successfully"})
}
