package handlers

import (
	"net/http"
	"time"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

const CurrentLicenseVersion = "v0.0.1"

// GetLicenseStatus handles GET /api/license/status
func GetLicenseStatus(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok || username == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized session"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User profile not found"})
	}

	// Non-admin users do not need to agree to administrative licenses
	if user.Role != "admin" {
		return c.JSON(http.StatusOK, echo.Map{
			"accepted":        true,
			"current_version": CurrentLicenseVersion,
		})
	}

	var agreement models.LicenseAgreement
	err := database.DB.Where("user_id = ? AND version = ?", user.ID, CurrentLicenseVersion).First(&agreement).Error
	if err != nil {
		// Not accepted yet
		return c.JSON(http.StatusOK, echo.Map{
			"accepted":        false,
			"current_version": CurrentLicenseVersion,
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"accepted":        true,
		"current_version": CurrentLicenseVersion,
		"accepted_at":     agreement.AcceptedAt,
	})
}

// AcceptLicense handles POST /api/license/accept
func AcceptLicense(c echo.Context) error {
	username, ok := c.Get("username").(string)
	if !ok || username == "" {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Unauthorized session"})
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "User profile not found"})
	}

	if user.Role != "admin" {
		return c.JSON(http.StatusForbidden, echo.Map{"error": "Only administrators can accept licenses"})
	}

	// Check if already accepted
	var existing models.LicenseAgreement
	err := database.DB.Where("user_id = ? AND version = ?", user.ID, CurrentLicenseVersion).First(&existing).Error
	if err == nil {
		return c.JSON(http.StatusOK, echo.Map{
			"message":  "License already accepted",
			"version":  CurrentLicenseVersion,
			"accepted": true,
		})
	}

	agreement := models.LicenseAgreement{
		UserID:     user.ID,
		Version:    CurrentLicenseVersion,
		AcceptedAt: time.Now().Unix(),
	}

	if err := database.DB.Create(&agreement).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to record license agreement: " + err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":  "License agreement accepted and registered successfully",
		"version":  CurrentLicenseVersion,
		"accepted": true,
	})
}
