package handlers

import (
	"log"
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

// RecordAuditLog inserts an audit trace record into the database.
func RecordAuditLog(assetID string, action string, changes string) {
	logRecord := models.AuditLog{
		AssetID: assetID,
		Action:  action,
		Changes: changes,
	}
	if err := database.DB.Create(&logRecord).Error; err != nil {
		log.Printf("[AuditLog Error] Failed to write audit trace: %v", err)
	}
}

// ReadAuditLogs handles GET /api/assets/{asset_id}/logs
func ReadAuditLogs(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Verify asset exists first
	var asset models.Asset
	if err := database.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	var logs []models.AuditLog
	if err := database.DB.Where("asset_id = ?", assetID).Order("timestamp desc").Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch audit logs"})
	}

	return c.JSON(http.StatusOK, logs)
}

// ReadGlobalAuditLogs handles GET /api/audit-logs
func ReadGlobalAuditLogs(c echo.Context) error {
	var logs []models.AuditLog
	if err := database.DB.Order("timestamp desc").Limit(100).Find(&logs).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch global audit logs"})
	}

	return c.JSON(http.StatusOK, logs)
}
