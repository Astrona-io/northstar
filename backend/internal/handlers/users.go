package handlers

import (
	"encoding/json"
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
)

// ReadUsers handles GET /api/users (Admin & Operator only)
func ReadUsers(c echo.Context) error {
	var users []models.User
	if err := database.DB.Preload("Group.Permissions").Preload("Permissions").Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch users: " + err.Error()})
	}
	return c.JSON(http.StatusOK, users)
}

// ReadGroups handles GET /api/groups (Admin & Operator only)
func ReadGroups(c echo.Context) error {
	var groups []models.Group
	if err := database.DB.Preload("Permissions").Find(&groups).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch groups: " + err.Error()})
	}
	return c.JSON(http.StatusOK, groups)
}

// ReadPermissions handles GET /api/permissions (Admin & Operator only)
func ReadPermissions(c echo.Context) error {
	var perms []models.Permission
	database.DB.Order("name asc, effect asc").Find(&perms)
	return c.JSON(http.StatusOK, perms)
}

type UserUpdateReq struct {
	GroupID       *uint  `json:"group_id"`
	PermissionIDs []uint `json:"permission_ids"` // User-level overrides!
}

// UpdateUser handles PUT /api/users/{id} (Admins only)
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid user ID"})
	}

	var user models.User
	if err := database.DB.Preload("Permissions").First(&user, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"error": "User not found"})
	}

	var body map[string]any
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON payload"})
	}

	updates := make(map[string]any)
	if val, exists := body["group_id"]; exists {
		if val == nil {
			updates["group_id"] = nil
		} else if s, ok := val.(string); ok {
			updates["group_id"] = &s
		}
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to update user: " + err.Error()})
		}
	}

	// Update many-to-many permission overrides (Phase 1 Admin User overrides)
	if permVal, exists := body["permission_ids"]; exists {
		if idsInterface, ok := permVal.([]any); ok {
			var permIDs []string
			for _, item := range idsInterface {
				if s, ok := item.(string); ok {
					permIDs = append(permIDs, s)
				}
			}

			var newPerms []models.Permission
			if len(permIDs) > 0 {
				database.DB.Find(&newPerms, permIDs)
			}
			// Replace GORM association
			database.DB.Model(&user).Association("Permissions").Replace(&newPerms)
		}
	}

	// Reload user with GORM Preloads
	database.DB.Preload("Group.Permissions").Preload("Permissions").First(&user, "id = ?", id)

	return c.JSON(http.StatusOK, user)
}
