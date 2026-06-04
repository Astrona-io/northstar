package handlers

import (
	"net/http"

	"cmdb-backend/internal/database"
	"cmdb-backend/internal/models"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// ReadCategories handles GET /api/categories (Phase 1 Dynamic Categories with Order/A-Z Sort)
func ReadCategories(c echo.Context) error {
	var categories []models.Category
	// Dual priority sort (Phase 1 Sorting): Custom weight "item_order" first, then Alphabetical "name"
	query := database.DB.Preload("SubGroups", func(db *gorm.DB) *gorm.DB {
		return db.Order("item_order asc, name asc")
	}).Order("item_order asc, name asc")

	if err := query.Find(&categories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to fetch categories: " + err.Error()})
	}
	return c.JSON(http.StatusOK, categories)
}

// CreateCategory handles POST /api/categories (Admin Only, Phase 1)
func CreateCategory(c echo.Context) error {
	var cat models.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if cat.Name == "" || cat.Icon == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Category Name and Icon are required"})
	}

	if err := database.DB.Create(&cat).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create category: " + err.Error()})
	}
	return c.JSON(http.StatusOK, cat)
}

// DeleteCategory handles DELETE /api/categories/:id (Admin Only)
func DeleteCategory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid ID"})
	}

	if err := database.DB.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to delete category: " + err.Error()})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": "Category deleted successfully"})
}

// CreateSubGroup handles POST /api/categories/sub-groups (Admin Only, Phase 1)
func CreateSubGroup(c echo.Context) error {
	var sub models.SubGroup
	if err := c.Bind(&sub); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid request body"})
	}
	if sub.CategoryID == "" || sub.Name == "" || sub.Icon == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Parent CategoryID, Name, and Icon are required"})
	}

	if err := database.DB.Create(&sub).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to create sub-group: " + err.Error()})
	}
	return c.JSON(http.StatusOK, sub)
}
