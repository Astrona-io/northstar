package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

// AssetCreateReq defines the payload structure for creating a new Asset.
type AssetCreateReq struct {
	Name                 string         `json:"name"`
	Type                 string         `json:"type"`
	Status               string         `json:"status"`
	IPAddress            *string        `json:"ip_address"`
	Description          *string        `json:"description"`
	Properties           models.JSONMap `json:"properties"`
	RackID               *string        `json:"rack_id"`
	RackPositionU        *int           `json:"rack_position_u"`
	HostAssetID          *string        `json:"host_asset_id"`
	ContainerImage       string         `json:"container_image"`
	ContainerPortMapping string         `json:"container_port_mapping"`
	DeviceModelID        *string        `json:"device_model_id"`
	DeviceModelRevision  *int           `json:"device_model_revision"`
}

// CreateAsset handles POST /api/assets/
func CreateAsset(c echo.Context) error {
	var req AssetCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name field is required"})
	}
	if req.Type == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Type field is required"})
	}
	if req.Status == "" {
		req.Status = "active"
	}
	if req.Properties == nil {
		req.Properties = make(models.JSONMap)
	}

	sanitizedProps, err := sanitizeProperties(req.Type, req.Properties)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid properties format"})
	}

	// Enforce Dynamic Category Dependency Chain Validations (Phase 1 Taxonomy Chains)
	var category models.Category
	if err := database.DB.Where("name = ?", req.Type).First(&category).Error; err == nil {
		if category.ParentDependency != "" {
			if req.HostAssetID == nil || *req.HostAssetID == "" {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation Error: This category type strictly requires an active parent worker/host of type " + category.ParentDependency})
			}
			var parent models.Asset
			if err := database.DB.Where("id = ?", *req.HostAssetID).First(&parent).Error; err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation Error: Referenced parent worker/host not found inside GORM"})
			}
			if parent.Type != category.ParentDependency {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation Error: This category type strictly requires a parent worker/host of type " + category.ParentDependency})
			}
		}
	}

	asset := models.Asset{
		Name:                 req.Name,
		Type:                 req.Type,
		Status:               req.Status,
		IPAddress:            req.IPAddress,
		Description:          req.Description,
		Properties:           sanitizedProps,
		RackID:               req.RackID,
		RackPositionU:        req.RackPositionU,
		HostAssetID:          req.HostAssetID,
		ContainerImage:       req.ContainerImage,
		ContainerPortMapping: req.ContainerPortMapping,
	}

	if req.DeviceModelID != nil && *req.DeviceModelID != "" {
		var dm models.DeviceModel
		if err := database.DB.First(&dm, "id = ?", *req.DeviceModelID).Error; err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{"error": "Referenced Device Model not found"})
		}
		asset.DeviceModelID = req.DeviceModelID
		if req.DeviceModelRevision != nil {
			asset.DeviceModelRevision = req.DeviceModelRevision
		} else {
			rev := dm.Revision
			asset.DeviceModelRevision = &rev
		}
	}

	if err := database.DB.Create(&asset).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create asset: " + err.Error()})
	}

	populateSingleAssetRevisionDetails(&asset)

	RecordAuditLog(asset.ID, "create", "Asset created with type: "+asset.Type)

	// Trigger Outbound Webhooks (Phase 3 Webhooks)
	DispatchWebhookEvent("asset:create", asset)

	asset.MaintenanceStatus = "none"
	return c.JSON(http.StatusOK, asset)
}

// ReadAssets handles GET /api/assets/ (with pagination, search, and custom filters)
func ReadAssets(c echo.Context) error {
	skipStr := c.QueryParam("skip")
	limitStr := c.QueryParam("limit")
	search := c.QueryParam("search")
	assetType := c.QueryParam("type")
	status := c.QueryParam("status")
	dcID := c.QueryParam("datacenter_id")
	mappingState := c.QueryParam("mapping_state") // all, mapped, orphaned (Phase 1)

	skip := 0
	limit := 100

	if skipStr != "" {
		if s, err := strconv.Atoi(skipStr); err == nil {
			skip = s
		}
	}
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	query := database.DB.Model(&models.Asset{}).Preload("MaintenanceWindows").Preload("Rack.Datacenter").Preload("DeviceModel.Manufacturer")

	if search != "" {
		searchTerm := "%" + search + "%"
		query = query.Where("name LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}
	if assetType != "" {
		query = query.Where("type = ?", assetType)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by Datacenter ID (via joining Racks, Phase 1 UI L3)
	if dcID != "" {
		query = query.Joins("LEFT JOIN racks ON racks.id = assets.rack_id").Where("racks.datacenter_id = ?", dcID)
	}

	// Filter by Mapping State (Orphan detection, Phase 1 UI L3)
	if mappingState == "mapped" {
		query = query.Where("rack_id IS NOT NULL OR host_asset_id IS NOT NULL")
	} else if mappingState == "orphaned" {
		query = query.Where("rack_id IS NULL AND host_asset_id IS NULL")
	}

	var assets []models.Asset
	err := query.Offset(skip).Limit(limit).Find(&assets).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch assets: " + err.Error()})
	}

	for i := range assets {
		assets[i].MaintenanceStatus = assets[i].ComputeMaintenanceStatus()
	}

	populateDeviceModelRevisionDetails(assets)

	return c.JSON(http.StatusOK, assets)
}

// ReadAsset handles GET /api/assets/{asset_id}
func ReadAsset(c echo.Context) error {
	id := c.Param("asset_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	var asset models.Asset
	err := database.DB.Preload("MaintenanceWindows").Preload("DeviceModel.Manufacturer").First(&asset, "id = ?", id).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	asset.MaintenanceStatus = asset.ComputeMaintenanceStatus()
	populateSingleAssetRevisionDetails(&asset)

	return c.JSON(http.StatusOK, asset)
}

func populateDeviceModelRevisionDetails(assets []models.Asset) {
	for i := range assets {
		if assets[i].DeviceModelID != nil && assets[i].DeviceModelRevision != nil {
			var rev models.DeviceModelRevision
			if err := database.DB.First(&rev, "device_model_id = ? AND revision = ?", *assets[i].DeviceModelID, *assets[i].DeviceModelRevision).Error; err == nil {
				assets[i].DeviceModelRevisionDetails = &rev
			}
		}
	}
}

func populateSingleAssetRevisionDetails(asset *models.Asset) {
	if asset.DeviceModelID != nil && asset.DeviceModelRevision != nil {
		var rev models.DeviceModelRevision
		if err := database.DB.First(&rev, "device_model_id = ? AND revision = ?", *asset.DeviceModelID, *asset.DeviceModelRevision).Error; err == nil {
			asset.DeviceModelRevisionDetails = &rev
		}
	}
}

// UpdateAsset handles PUT /api/assets/{asset_id}
// Implements partial updates based strictly on fields provided in the JSON body.
func UpdateAsset(c echo.Context) error {
	id := c.Param("asset_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	var asset models.Asset
	err := database.DB.Preload("MaintenanceWindows").First(&asset, "id = ?", id).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	// Bind raw map to identify only supplied fields for partial updates
	var body map[string]any
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON: " + err.Error()})
	}

	updates := make(map[string]any)
	if val, exists := body["name"]; exists {
		if s, ok := val.(string); ok {
			updates["name"] = s
		}
	}
	if val, exists := body["type"]; exists {
		if s, ok := val.(string); ok {
			updates["type"] = s
		}
	}
	if val, exists := body["status"]; exists {
		if s, ok := val.(string); ok {
			updates["status"] = s
		}
	}
	if val, exists := body["ip_address"]; exists {
		if val == nil {
			updates["ip_address"] = nil
		} else if s, ok := val.(string); ok {
			updates["ip_address"] = s
		}
	}
	if val, exists := body["description"]; exists {
		if val == nil {
			updates["description"] = nil
		} else if s, ok := val.(string); ok {
			updates["description"] = s
		}
	}
	if val, exists := body["rack_id"]; exists {
		if val == nil {
			updates["rack_id"] = nil
		} else if s, ok := val.(string); ok {
			updates["rack_id"] = s
		}
	}
	if val, exists := body["rack_position_u"]; exists {
		if val == nil {
			updates["rack_position_u"] = nil
		} else {
			switch v := val.(type) {
			case float64:
				ival := int(v)
				updates["rack_position_u"] = &ival
			case int:
				updates["rack_position_u"] = &v
			}
		}
	}
	if val, exists := body["host_asset_id"]; exists {
		if val == nil {
			updates["host_asset_id"] = nil
		} else if s, ok := val.(string); ok {
			updates["host_asset_id"] = s
		}
	}
	if val, exists := body["container_image"]; exists {
		if s, ok := val.(string); ok {
			updates["container_image"] = s
		}
	}
	if val, exists := body["container_port_mapping"]; exists {
		if s, ok := val.(string); ok {
			updates["container_port_mapping"] = s
		}
	}
	if val, exists := body["device_model_id"]; exists {
		if val == nil {
			updates["device_model_id"] = nil
			updates["device_model_revision"] = nil
		} else if s, ok := val.(string); ok {
			updates["device_model_id"] = s
			if _, revExists := body["device_model_revision"]; !revExists {
				var dm models.DeviceModel
				if err := database.DB.First(&dm, "id = ?", s).Error; err == nil {
					updates["device_model_revision"] = dm.Revision
				}
			}
		}
	}
	if val, exists := body["device_model_revision"]; exists {
		if val == nil {
			updates["device_model_revision"] = nil
		} else {
			switch v := val.(type) {
			case float64:
				ival := int(v)
				updates["device_model_revision"] = ival
			case int:
				updates["device_model_revision"] = v
			}
		}
	}
	if val, exists := body["properties"]; exists {
		if val == nil {
			updates["properties"] = models.JSONMap{}
		} else if m, ok := val.(map[string]any); ok {
			assetType := asset.Type
			if typeVal, typeExists := updates["type"]; typeExists {
				assetType = typeVal.(string)
			}
			sanitizedProps, err := sanitizeProperties(assetType, models.JSONMap(m))
			if err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid properties format"})
			}
			updates["properties"] = sanitizedProps
		}
	}

	// Enforce Dynamic Category Dependency Chain Validations (Phase 1 Taxonomy Chains)
	prospectiveType := asset.Type
	if tVal, exists := updates["type"]; exists {
		prospectiveType = tVal.(string)
	}

	var category models.Category
	if err := database.DB.Where("name = ?", prospectiveType).First(&category).Error; err == nil {
		if category.ParentDependency != "" {
			// Find prospective Host ID
			var prospectiveHostID *string = asset.HostAssetID
			if hVal, exists := updates["host_asset_id"]; exists {
				if hVal == nil {
					prospectiveHostID = nil
				} else {
					str := hVal.(string)
					prospectiveHostID = &str
				}
			}

			if prospectiveHostID == nil || *prospectiveHostID == "" {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation Error: This category type strictly requires an active parent worker/host of type " + category.ParentDependency})
			}

			var parent models.Asset
			if err := database.DB.Where("id = ?", *prospectiveHostID).First(&parent).Error; err != nil {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation Error: Referenced parent worker/host not found inside GORM"})
			}
			if parent.Type != category.ParentDependency {
				return c.JSON(http.StatusBadRequest, echo.Map{"error": "Validation Error: This category type strictly requires a parent worker/host of type " + category.ParentDependency})
			}
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&asset).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update asset: " + err.Error()})
		}

		// Serialize changes for audit trail (Phase 2)
		b, _ := json.Marshal(updates)
		RecordAuditLog(asset.ID, "update", string(b))
	}

	// Reload to get fully updated fields + preload windows for status recalculation
	if err := database.DB.Preload("MaintenanceWindows").Preload("DeviceModel.Manufacturer").First(&asset, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to reload updated asset"})
	}
	asset.MaintenanceStatus = asset.ComputeMaintenanceStatus()
	populateSingleAssetRevisionDetails(&asset)

	// Trigger Outbound Webhooks (Phase 3 Webhooks)
	DispatchWebhookEvent("asset:update", asset)

	return c.JSON(http.StatusOK, asset)
}

// DeleteAsset handles DELETE /api/assets/{asset_id}
func DeleteAsset(c echo.Context) error {
	id := c.Param("asset_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	var asset models.Asset
	err := database.DB.First(&asset, "id = ?", id).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	if err := database.DB.Delete(&asset).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete asset: " + err.Error()})
	}

	RecordAuditLog(asset.ID, "delete", "Asset deleted: "+asset.Name)

	// Trigger Outbound Webhooks (Phase 3 Webhooks)
	DispatchWebhookEvent("asset:delete", asset)

	return c.JSON(http.StatusOK, echo.Map{"message": "Asset deleted successfully"})
}
