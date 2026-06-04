package handlers

import (
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type RelationshipSyncReq struct {
	UplinkIDs   []string `json:"uplink_ids"`
	DownlinkIDs []string `json:"downlink_ids"`
}

// SyncRelationships handles POST /api/assets/{asset_id}/relationships
func SyncRelationships(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Verify asset exists first
	var asset models.Asset
	if err := database.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	var req RelationshipSyncReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}

	// Delete existing relationships for this asset
	database.DB.Where("source_id = ?", assetID).Delete(&models.AssetRelationship{})

	// Save Uplinks (Source has uplink to Target)
	for _, upID := range req.UplinkIDs {
		if upID == assetID {
			continue // Prevent self-linking
		}
		rel := models.AssetRelationship{
			SourceID: assetID,
			TargetID: upID,
			Type:     "uplink",
		}
		database.DB.Create(&rel)
	}

	// Save Downlinks (Source has downlink to Target)
	for _, downID := range req.DownlinkIDs {
		if downID == assetID {
			continue // Prevent self-linking
		}
		rel := models.AssetRelationship{
			SourceID: assetID,
			TargetID: downID,
			Type:     "downlink",
		}
		database.DB.Create(&rel)
	}

	// Log audit trace (Phase 2)
	RecordAuditLog(assetID, "update", "Synchronized network topology connections")

	return c.JSON(http.StatusOK, echo.Map{"message": "Relationships synchronized successfully"})
}

// ReadRelationships handles GET /api/assets/{asset_id}/relationships
func ReadRelationships(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Find uplinks (assets that this asset has uplink to)
	var uplinkRels []models.AssetRelationship
	database.DB.Where("source_id = ? AND type = 'uplink'", assetID).Find(&uplinkRels)

	var uplinkIDs []string
	for _, r := range uplinkRels {
		uplinkIDs = append(uplinkIDs, r.TargetID)
	}

	var uplinks []models.Asset
	if len(uplinkIDs) > 0 {
		database.DB.Find(&uplinks, uplinkIDs)
	}

	// Find downlinks (assets that this asset has downlink to)
	var downlinkRels []models.AssetRelationship
	database.DB.Where("source_id = ? AND type = 'downlink'", assetID).Find(&downlinkRels)

	var downlinkIDs []string
	for _, r := range downlinkRels {
		downlinkIDs = append(downlinkIDs, r.TargetID)
	}

	var downlinks []models.Asset
	if len(downlinkIDs) > 0 {
		database.DB.Find(&downlinks, downlinkIDs)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"uplinks":   uplinks,
		"downlinks": downlinks,
	})
}
