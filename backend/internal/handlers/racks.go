package handlers

import (
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

type RackCoordinatesReq struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// UpdateRackCoordinates handles PUT /api/racks/:id/coordinates (Phase 1 CAD Designer)
func UpdateRackCoordinates(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid rack ID"})
	}

	var req RackCoordinatesReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid coordinates payload"})
	}

	var rack models.Rack
	if err := database.DB.First(&rack, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "Rack not found"})
	}

	// Update GORM X/Y bounds
	rack.X = req.X
	rack.Y = req.Y

	if err := database.DB.Save(&rack).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update rack coordinates: " + err.Error()})
	}

	return c.JSON(http.StatusOK, rack)
}
