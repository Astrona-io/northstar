package handlers

import (
	"encoding/json"
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type ManufacturerCreateReq struct {
	Name    string `json:"name"`
	Website string `json:"website"`
}

// CreateManufacturer handles POST /api/manufacturers
func CreateManufacturer(c echo.Context) error {
	var req ManufacturerCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name field is required"})
	}

	m := models.Manufacturer{Name: req.Name, Website: req.Website}
	if err := database.DB.Create(&m).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Manufacturer brand already exists or failed to save"})
	}
	return c.JSON(http.StatusOK, m)
}

// ReadManufacturers handles GET /api/manufacturers
func ReadManufacturers(c echo.Context) error {
	var list []models.Manufacturer
	database.DB.Order("name asc").Find(&list)
	return c.JSON(http.StatusOK, list)
}

type DeviceModelCreateReq struct {
	ManufacturerID string         `json:"manufacturer_id"`
	ModelName      string         `json:"model_name"`
	CategoryIDs    []string       `json:"category_ids"`
	GeneralInfo    string         `json:"general_info"`
	Ports          models.JSONMap `json:"ports"`
}

// CreateDeviceModel handles POST /api/devices/
func CreateDeviceModel(c echo.Context) error {
	var req DeviceModelCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}

	if req.ManufacturerID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ManufacturerID field is required"})
	}
	if req.ModelName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "ModelName field is required"})
	}

	// Verify Manufacturer exists first (Strict Manufacturer Profile - Phase 2 Catalog)
	var m models.Manufacturer
	if err := database.DB.First(&m, "id = ?", req.ManufacturerID).Error; err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Selected Manufacturer brand does not exist"})
	}

	// Resolve dynamic Categories tags (Multi-Category Classification - Phase 2 Catalog)
	var cats []models.Category
	if len(req.CategoryIDs) > 0 {
		database.DB.Find(&cats, req.CategoryIDs)
	} else {
		// Fallback/Default to "global" category (Phase 2 Catalog)
		var globalCat models.Category
		database.DB.Where("name = ?", "global").First(&globalCat)
		cats = append(cats, globalCat)
	}

	device := models.DeviceModel{
		ManufacturerID: req.ManufacturerID,
		ModelName:      req.ModelName,
		Categories:     cats,
		GeneralInfo:    req.GeneralInfo,
		Revision:       1,
		Ports:          req.Ports,
	}

	if err := database.DB.Create(&device).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create device spec: " + err.Error()})
	}

	// Create first revision record
	if err := saveOrUpdateRevisionRecord(database.DB, &device); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create first device revision: " + err.Error()})
	}

	// Reload to fetch Preloaded associations
	database.DB.Preload("Manufacturer").Preload("Categories").First(&device, "id = ?", device.ID)

	return c.JSON(http.StatusOK, device)
}

func saveOrUpdateRevisionRecord(db *gorm.DB, dm *models.DeviceModel) error {
	var rev models.DeviceModelRevision
	err := db.First(&rev, "device_model_id = ? AND revision = ?", dm.ID, dm.Revision).Error
	if err == nil {
		rev.ModelName = dm.ModelName
		rev.GeneralInfo = dm.GeneralInfo
		rev.Ports = dm.Ports
		return db.Save(&rev).Error
	}
	rev = models.DeviceModelRevision{
		DeviceModelID: dm.ID,
		Revision:      dm.Revision,
		ModelName:     dm.ModelName,
		GeneralInfo:   dm.GeneralInfo,
		Ports:         dm.Ports,
	}
	return db.Create(&rev).Error
}

// ReadDeviceModels handles GET /api/devices/
func ReadDeviceModels(c echo.Context) error {
	var devices []models.DeviceModel
	if err := database.DB.Preload("Manufacturer").Preload("Categories").Find(&devices).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch devices: " + err.Error()})
	}

	return c.JSON(http.StatusOK, devices)
}

// ReadDeviceModel handles GET /api/devices/{device_id}
func ReadDeviceModel(c echo.Context) error {
	id := c.Param("device_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid device ID"})
	}

	var device models.DeviceModel
	if err := database.DB.Preload("Manufacturer").Preload("Categories").First(&device, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Device spec not found"})
	}

	return c.JSON(http.StatusOK, device)
}

// UpdateDeviceModel handles PUT /api/devices/{device_id}
func UpdateDeviceModel(c echo.Context) error {
	id := c.Param("device_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid device ID"})
	}

	var device models.DeviceModel
	if err := database.DB.Preload("Categories").First(&device, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Device not found"})
	}

	var body map[string]any
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON: " + err.Error()})
	}

	hasChanges := false

	newManufacturerID := device.ManufacturerID
	if val, exists := body["manufacturer_id"]; exists {
		if s, ok := val.(string); ok && s != device.ManufacturerID {
			newManufacturerID = s
			hasChanges = true
		}
	}

	newModelName := device.ModelName
	if val, exists := body["model_name"]; exists {
		if s, ok := val.(string); ok && s != device.ModelName {
			newModelName = s
			hasChanges = true
		}
	}

	newGeneralInfo := device.GeneralInfo
	if val, exists := body["general_info"]; exists {
		if s, ok := val.(string); ok && s != device.GeneralInfo {
			newGeneralInfo = s
			hasChanges = true
		}
	}

	newPorts := device.Ports
	if val, exists := body["ports"]; exists {
		portsMap := make(models.JSONMap)
		if m, ok := val.(map[string]any); ok {
			for k, v := range m {
				portsMap[k] = v
			}
		}
		oldBytes, _ := json.Marshal(device.Ports)
		newBytes, _ := json.Marshal(portsMap)
		if string(oldBytes) != string(newBytes) {
			newPorts = portsMap
			hasChanges = true
		}
	}

	if hasChanges {
		var assetCount int64
		database.DB.Model(&models.Asset{}).Where("device_model_id = ?", id).Count(&assetCount)

		if assetCount > 0 {
			// Bump revision number!
			device.Revision = device.Revision + 1
		}

		device.ManufacturerID = newManufacturerID
		device.ModelName = newModelName
		device.GeneralInfo = newGeneralInfo
		device.Ports = newPorts

		if err := database.DB.Save(&device).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update device: " + err.Error()})
		}

		if err := saveOrUpdateRevisionRecord(database.DB, &device); err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to save device revision: " + err.Error()})
		}
	}

	// Update Categories many-to-many associations if provided in body
	if catVal, exists := body["category_ids"]; exists {
		if idsInterface, ok := catVal.([]any); ok {
			var catIDs []string
			for _, item := range idsInterface {
				if s, ok := item.(string); ok {
					catIDs = append(catIDs, s)
				}
			}

			var newCats []models.Category
			if len(catIDs) > 0 {
				database.DB.Find(&newCats, catIDs)
			} else {
				var globalCat models.Category
				database.DB.Where("name = ?", "global").First(&globalCat)
				newCats = append(newCats, globalCat)
			}
			// Replace association in DB (clean many-to-many update)
			database.DB.Model(&device).Association("Categories").Replace(&newCats)
		}
	}

	if err := database.DB.Preload("Manufacturer").Preload("Categories").First(&device, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to reload device"})
	}

	return c.JSON(http.StatusOK, device)
}

// DeleteDeviceModel handles DELETE /api/devices/{device_id}
func DeleteDeviceModel(c echo.Context) error {
	id := c.Param("device_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid device ID"})
	}

	var device models.DeviceModel
	if err := database.DB.First(&device, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Device not found"})
	}

	// Clear many-to-many categories association first
	database.DB.Model(&device).Association("Categories").Clear()

	if err := database.DB.Delete(&device).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete device: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Device deleted successfully"})
}
