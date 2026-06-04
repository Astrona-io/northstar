package handlers

import (
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

// ReadDatacenterTypes handles GET /api/datacenter-types (Phase 3)
func ReadDatacenterTypes(c echo.Context) error {
	var types []models.DatacenterType
	if err := database.DB.Order("name asc").Find(&types).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch datacenter types"})
	}
	return c.JSON(http.StatusOK, types)
}

// CreateDatacenterType handles POST /api/datacenter-types (Phase 3)
func CreateDatacenterType(c echo.Context) error {
	var dt models.DatacenterType
	if err := c.Bind(&dt); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if dt.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Datacenter type name is required"})
	}
	if err := database.DB.Create(&dt).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create datacenter type"})
	}
	return c.JSON(http.StatusOK, dt)
}

// ReadDatacenters handles GET /api/datacenters
func ReadDatacenters(c echo.Context) error {
	var datacenters []models.Datacenter
	// Preload DatacenterType, Floors, and nested Racks/Walls (DCIM CAD L4)
	if err := database.DB.Preload("DatacenterType").Preload("Floors.Racks").Preload("Floors.Walls").Preload("Racks").Find(&datacenters).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch datacenters: " + err.Error()})
	}
	return c.JSON(http.StatusOK, datacenters)
}

// CreateDatacenter handles POST /api/datacenters (Admin Only)
func CreateDatacenter(c echo.Context) error {
	var dc models.Datacenter
	if err := c.Bind(&dc); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if dc.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Name is a required field"})
	}

	if err := database.DB.Create(&dc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create datacenter or name already exists: " + err.Error()})
	}

	return c.JSON(http.StatusOK, dc)
}

// CreateDatacenterFloor handles POST /api/datacenter-floors (Admin Only)
func CreateDatacenterFloor(c echo.Context) error {
	var floor models.DatacenterFloor
	if err := c.Bind(&floor); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if floor.DatacenterID == "" || floor.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Datacenter ID and Floor Name are required"})
	}
	// Defaults for CAD Floor Planner
	if floor.Width == 0 {
		floor.Width = 1000 // default 1000cm / 10m
	}
	if floor.Depth == 0 {
		floor.Depth = 1000 // default 1000cm / 10m
	}

	if err := database.DB.Create(&floor).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create datacenter floor"})
	}
	return c.JSON(http.StatusOK, floor)
}

// UpdateDatacenterFloor handles PUT /api/datacenter-floors/:id (Admin Only)
func UpdateDatacenterFloor(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid floor ID"})
	}

	var existing models.DatacenterFloor
	if err := database.DB.First(&existing, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Floor not found"})
	}

	var update struct {
		Name  *string  `json:"name"`
		Level *int     `json:"level"`
		Width *float64 `json:"width"`
		Depth *float64 `json:"depth"`
	}
	if err := c.Bind(&update); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if update.Name != nil {
		existing.Name = *update.Name
	}
	if update.Level != nil {
		existing.Level = *update.Level
	}
	if update.Width != nil {
		existing.Width = *update.Width
	}
	if update.Depth != nil {
		existing.Depth = *update.Depth
	}

	if err := database.DB.Save(&existing).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update floor: " + err.Error()})
	}

	return c.JSON(http.StatusOK, existing)
}

// DeleteDatacenterFloor handles DELETE /api/datacenter-floors/:id (Admin Only)
func DeleteDatacenterFloor(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid floor ID"})
	}
	if err := database.DB.Delete(&models.DatacenterFloor{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete datacenter floor"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Datacenter floor deleted successfully"})
}

// DeleteDatacenter handles DELETE /api/datacenters/:id (Admin Only)
func DeleteDatacenter(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid datacenter ID"})
	}

	var dc models.Datacenter
	if err := database.DB.First(&dc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Datacenter not found"})
	}

	if err := database.DB.Delete(&dc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete datacenter: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Datacenter deleted successfully"})
}

// ReadRacks handles GET /api/racks
func ReadRacks(c echo.Context) error {
	var racks []models.Rack
	if err := database.DB.Preload("Datacenter").Find(&racks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch racks: " + err.Error()})
	}
	return c.JSON(http.StatusOK, racks)
}

// CreateRack handles POST /api/racks (Admin Only)
func CreateRack(c echo.Context) error {
	var rack models.Rack
	if err := c.Bind(&rack); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}

	if rack.Name == "" || rack.DatacenterID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Rack Name and Datacenter ID are required"})
	}

	if rack.HeightU <= 0 {
		rack.HeightU = 42
	}

	if err := database.DB.Create(&rack).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create cabinet: " + err.Error()})
	}

	return c.JSON(http.StatusOK, rack)
}

// DeleteRack handles DELETE /api/racks/:id (Admin Only)
func DeleteRack(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid rack ID"})
	}

	var rack models.Rack
	if err := database.DB.First(&rack, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Rack cabinet not found"})
	}

	if err := database.DB.Delete(&rack).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete rack cabinet: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": "Rack cabinet deleted successfully"})
}

// ReadRack handles GET /api/racks/:id
func ReadRack(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid rack ID"})
	}

	var rack models.Rack
	// Preload parent Datacenter and all nested assets mounted inside this cabinet (Phase 6 DCIM Mapping)
	if err := database.DB.Preload("Datacenter").Preload("Assets").First(&rack, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Rack cabinet not found"})
	}

	return c.JSON(http.StatusOK, rack)
}

// CreateDatacenterWall handles POST /api/dcim/walls (Admin Only, DCIM CAD L4)
func CreateDatacenterWall(c echo.Context) error {
	var wall models.DatacenterWall
	if err := c.Bind(&wall); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if wall.DatacenterFloorID == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "DatacenterFloorID is required"})
	}
	if err := database.DB.Create(&wall).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to save datacenter wall"})
	}
	return c.JSON(http.StatusOK, wall)
}

// DeleteDatacenterWall handles DELETE /api/dcim/walls/:id (Admin Only, DCIM CAD L4)
func DeleteDatacenterWall(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid wall ID"})
	}
	if err := database.DB.Delete(&models.DatacenterWall{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete datacenter wall"})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Wall deleted successfully"})
}
