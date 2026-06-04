package handlers

import (
	"encoding/json"
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type NetworkInterfaceCreateReq struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	MACAddress string `json:"mac_address"`
	IPAddress  string `json:"ip_address"`
	Speed      string `json:"speed"`
	MTU        int    `json:"mtu"`
	VLAN       string `json:"vlan"`
	Status     string `json:"status"`
}

// CreateNetworkInterface handles POST /api/assets/{asset_id}/interfaces
func CreateNetworkInterface(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Verify asset exists first
	var asset models.Asset
	if err := database.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	var req NetworkInterfaceCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name field is required"})
	}
	if req.MTU == 0 {
		req.MTU = 1500
	}
	if req.Status == "" {
		req.Status = "up"
	}

	netInterface := models.NetworkInterface{
		AssetID:    assetID,
		Name:       req.Name,
		Type:       req.Type,
		MACAddress: req.MACAddress,
		IPAddress:  req.IPAddress,
		Speed:      req.Speed,
		MTU:        req.MTU,
		VLAN:       req.VLAN,
		Status:     req.Status,
	}

	if err := database.DB.Create(&netInterface).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create network interface: " + err.Error()})
	}

	return c.JSON(http.StatusOK, netInterface)
}

// ReadNetworkInterfaces handles GET /api/assets/{asset_id}/interfaces
func ReadNetworkInterfaces(c echo.Context) error {
	assetID := c.Param("asset_id")
	if assetID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid asset ID"})
	}

	// Verify asset exists first
	var asset models.Asset
	if err := database.DB.First(&asset, "id = ?", assetID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Asset not found"})
	}

	var interfaces []models.NetworkInterface
	if err := database.DB.Where("asset_id = ?", assetID).Find(&interfaces).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch interfaces: " + err.Error()})
	}

	return c.JSON(http.StatusOK, interfaces)
}

// UpdateNetworkInterface handles PUT /api/interfaces/{interface_id}
func UpdateNetworkInterface(c echo.Context) error {
	id := c.Param("interface_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid interface ID"})
	}

	var netInterface models.NetworkInterface
	if err := database.DB.First(&netInterface, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Network interface not found"})
	}

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
	if val, exists := body["mac_address"]; exists {
		if s, ok := val.(string); ok {
			updates["mac_address"] = s
		}
	}
	if val, exists := body["ip_address"]; exists {
		if s, ok := val.(string); ok {
			updates["ip_address"] = s
		}
	}
	if val, exists := body["speed"]; exists {
		if s, ok := val.(string); ok {
			updates["speed"] = s
		}
	}
	if val, exists := body["mtu"]; exists {
		if n, ok := val.(float64); ok {
			updates["mtu"] = int(n)
		}
	}
	if val, exists := body["vlan"]; exists {
		if s, ok := val.(string); ok {
			updates["vlan"] = s
		}
	}
	if val, exists := body["status"]; exists {
		if s, ok := val.(string); ok {
			updates["status"] = s
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&netInterface).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update interface: " + err.Error()})
		}
	}

	if err := database.DB.First(&netInterface, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to reload updated interface"})
	}

	return c.JSON(http.StatusOK, netInterface)
}

// DeleteNetworkInterface handles DELETE /api/interfaces/{interface_id}
func DeleteNetworkInterface(c echo.Context) error {
	id := c.Param("interface_id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid interface ID"})
	}

	var netInterface models.NetworkInterface
	if err := database.DB.First(&netInterface, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Network interface not found"})
	}

	if err := database.DB.Delete(&netInterface).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete interface: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Interface deleted successfully"})
}
