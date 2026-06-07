package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type NetworkInterfaceCreateReq struct {
	NICName     string `json:"nic_name"`
	Type        string `json:"type"`
	MACAddress  string `json:"mac_address"`
	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`
	Speed       string `json:"speed"`
	MTU         int    `json:"mtu"`
	VLAN        string `json:"vlan"`
	Status      string `json:"status"`
	PortCount   int    `json:"port_count"`
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

	if req.NICName == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "NIC Name field is required"})
	}
	if req.MTU == 0 {
		req.MTU = 1500
	}
	if req.Status == "" {
		req.Status = "up"
	}
	if req.PortCount <= 0 {
		req.PortCount = 1
	}

	var createdPorts []models.NetworkInterface

	// Simple helper to format MAC increments
	incrementMAC := func(mac string, offset int) string {
		if len(mac) < 17 {
			return mac
		}
		lastByteHex := mac[15:17]
		var val int
		_, err := fmt.Sscanf(lastByteHex, "%X", &val)
		if err == nil {
			newVal := (val + offset) % 256
			return mac[:15] + fmt.Sprintf("%02X", newVal)
		}
		return mac
	}

	for i := 1; i <= req.PortCount; i++ {
		portName := req.NICName
		if req.PortCount > 1 {
			portName = req.NICName + " - Port " + strconv.Itoa(i)
		}
		macAddr := req.MACAddress
		if macAddr != "" && req.PortCount > 1 {
			macAddr = incrementMAC(macAddr, i-1)
		}

		// Only set IPAddress on Port 1 if multiple ports are created, to avoid duplicate IPs
		ipv4Addr := req.IPv4Address
		ipv6Addr := req.IPv6Address
		if i > 1 {
			ipv4Addr = ""
			ipv6Addr = ""
		}

		netInterface := models.NetworkInterface{
			AssetID:     assetID,
			NICName:     req.NICName,
			Name:        portName,
			Type:        req.Type,
			MACAddress:  macAddr,
			IPv4Address: ipv4Addr,
			IPv6Address: ipv6Addr,
			Speed:       req.Speed,
			MTU:         req.MTU,
			VLAN:        req.VLAN,
			Status:      req.Status,
		}

		if err := database.DB.Create(&netInterface).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create network interface port: " + err.Error()})
		}
		createdPorts = append(createdPorts, netInterface)
	}

	if len(createdPorts) > 0 {
		return c.JSON(http.StatusOK, createdPorts[0])
	}
	return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create any ports"})
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
	if val, exists := body["ipv4_address"]; exists {
		if s, ok := val.(string); ok {
			updates["ipv4_address"] = s
		}
	}
	if val, exists := body["ipv6_address"]; exists {
		if s, ok := val.(string); ok {
			updates["ipv6_address"] = s
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
	if val, exists := body["port_count"]; exists {
		if n, ok := val.(float64); ok {
			updates["port_count"] = int(n)
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

// ReadPortTypeProfiles handles GET /api/port-types
func ReadPortTypeProfiles(c echo.Context) error {
	var profiles []models.PortTypeProfile
	if err := database.DB.Order("type asc").Find(&profiles).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch port type profiles: " + err.Error()})
	}
	return c.JSON(http.StatusOK, profiles)
}

// CreatePortTypeProfile handles POST /api/port-types
func CreatePortTypeProfile(c echo.Context) error {
	var req models.PortTypeProfile
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body: " + err.Error()})
	}
	if req.Type == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "type field is required"})
	}
	if req.Name == "" {
		req.Name = req.Type
	}

	if err := database.DB.Create(&req).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create profile: " + err.Error()})
	}
	return c.JSON(http.StatusOK, req)
}

// UpdatePortTypeProfile handles PUT /api/port-types/:id
func UpdatePortTypeProfile(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid profile ID"})
	}

	var existing models.PortTypeProfile
	if err := database.DB.First(&existing, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Profile not found"})
	}

	var body map[string]any
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON: " + err.Error()})
	}

	updates := make(map[string]any)
	if val, exists := body["type"]; exists {
		if s, ok := val.(string); ok {
			updates["type"] = s
		}
	}
	if val, exists := body["name"]; exists {
		if s, ok := val.(string); ok {
			updates["name"] = s
		}
	}
	if val, exists := body["speeds"]; exists {
		if s, ok := val.(string); ok {
			updates["speeds"] = s
		} else if arr, ok := val.([]any); ok {
			speedsJson, _ := json.Marshal(arr)
			updates["speeds"] = string(speedsJson)
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&existing).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update profile: " + err.Error()})
		}
	}

	database.DB.First(&existing, "id = ?", id)
	return c.JSON(http.StatusOK, existing)
}

// DeletePortTypeProfile handles DELETE /api/port-types/:id
func DeletePortTypeProfile(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid profile ID"})
	}

	var profile models.PortTypeProfile
	if err := database.DB.First(&profile, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Profile not found"})
	}

	if err := database.DB.Delete(&profile).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete profile: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Profile deleted successfully"})
}
