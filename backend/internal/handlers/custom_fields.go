package handlers

import (
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type CustomFieldCreateReq struct {
	Name            string `json:"name"`
	FieldType       string `json:"field_type"` // text, number, boolean
	AssetType       string `json:"asset_type"` // Server, Switch, etc.
	IsRequired      bool   `json:"is_required"`
	TabGroup        string `json:"tab_group"`
	ValidationRegex string `json:"validation_regex"`
}

// CreateCustomField handles POST /api/custom-fields (Admins only)
func CreateCustomField(c echo.Context) error {
	var req CustomFieldCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if req.Name == "" || req.FieldType == "" || req.AssetType == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name, FieldType, and AssetType fields are required"})
	}

	field := models.CustomField{
		Name:            req.Name,
		FieldType:       req.FieldType,
		AssetType:       req.AssetType,
		IsRequired:      req.IsRequired,
		TabGroup:        req.TabGroup,
		ValidationRegex: req.ValidationRegex,
	}

	if err := database.DB.Create(&field).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Custom field name already exists or failed to save"})
	}

	return c.JSON(http.StatusOK, field)
}

// ReadCustomFields handles GET /api/custom-fields
func ReadCustomFields(c echo.Context) error {
	assetType := c.QueryParam("asset_type")

	query := database.DB.Model(&models.CustomField{})
	if assetType != "" {
		query = query.Where("asset_type = ?", assetType)
	}

	var list []models.CustomField
	if err := query.Order("name asc").Find(&list).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch custom fields: " + err.Error()})
	}

	return c.JSON(http.StatusOK, list)
}

// DeleteCustomField handles DELETE /api/custom-fields/{id}
func DeleteCustomField(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid field ID"})
	}

	var field models.CustomField
	if err := database.DB.First(&field, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Custom field not found"})
	}

	if err := database.DB.Delete(&field).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete custom field: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Custom field deleted successfully"})
}
