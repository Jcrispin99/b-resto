package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetInventoryCategories godoc
// @Summary      Listar categorías de inventario
// @Description  Obtiene lista de todas las categorías de inventario (para materias primas)
// @Tags         inventory-categories
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        type       query  string  false  "Filtrar por tipo"
// @Success      200  {object}  map[string]interface{}  "data: array de inventory categories"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /inventory-categories [get]
// @Security     Bearer
func GetInventoryCategories(c *gin.Context) {
	var categories []models.InventoryCategory

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if catType := c.Query("type"); catType != "" {
		query = query.Where("type = ?", catType)
	}

	// Cargar jerarquía
	query = query.Preload("Parent").Preload("Children")

	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventory categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetInventoryCategory godoc
// @Summary      Obtener categoría de inventario
// @Description  Obtiene una categoría de inventario por ID con su jerarquía
// @Tags         inventory-categories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la categoría"
// @Success      200  {object}  map[string]interface{}  "data: inventory category"
// @Failure      404  {object}  map[string]string       "error: Inventory category not found"
// @Router       /inventory-categories/{id} [get]
// @Security     Bearer
func GetInventoryCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.InventoryCategory

	if err := config.DB.Preload("Parent").Preload("Children").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateInventoryCategory godoc
// @Summary      Crear categoría de inventario
// @Description  Crea una nueva categoría de inventario
// @Tags         inventory-categories
// @Accept       json
// @Produce      json
// @Param        category  body  models.InventoryCategory  true  "Datos de la categoría"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /inventory-categories [post]
// @Security     Bearer
func CreateInventoryCategory(c *gin.Context) {
	var category models.InventoryCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar parent_id si se proporciona
	if category.ParentID != nil {
		var parent models.InventoryCategory
		if err := config.DB.First(&parent, *category.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent category not found"})
			return
		}
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Inventory category created successfully",
		"data":    category,
	})
}

// UpdateInventoryCategory godoc
// @Summary      Actualizar categoría de inventario
// @Description  Actualiza los datos de una categoría de inventario existente
// @Tags         inventory-categories
// @Accept       json
// @Produce      json
// @Param        id        path  int                        true  "ID de la categoría"
// @Param        category  body  models.InventoryCategory   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Inventory category not found"
// @Router       /inventory-categories/{id} [put]
// @Security     Bearer
func UpdateInventoryCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.InventoryCategory

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory category not found"})
		return
	}

	var updateData models.InventoryCategory
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar que no se establezca como su propio padre
	if updateData.ParentID != nil && *updateData.ParentID == category.ID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category cannot be its own parent"})
		return
	}

	if err := config.DB.Model(&category).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory category updated successfully",
		"data":    category,
	})
}

// DeleteInventoryCategory godoc
// @Summary      Eliminar categoría de inventario
// @Description  Elimina una categoría de inventario (soft delete). No permite si tiene subcategorías.
// @Tags         inventory-categories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la categoría"
// @Success      200  {object}  map[string]string  "message: Inventory category deleted successfully"
// @Failure      400  {object}  map[string]string  "error: tiene subcategorías"
// @Failure      404  {object}  map[string]string  "error: Inventory category not found"
// @Router       /inventory-categories/{id} [delete]
// @Security     Bearer
func DeleteInventoryCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.InventoryCategory

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory category not found"})
		return
	}

	// Verificar si tiene subcategorías
	var childCount int64
	config.DB.Model(&models.InventoryCategory{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with subcategories"})
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete inventory category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inventory category deleted successfully"})
}

// ToggleInventoryCategoryStatus godoc
// @Summary      Activar/Desactivar categoría de inventario
// @Description  Cambia el estado is_active de una categoría de inventario
// @Tags         inventory-categories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la categoría"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Inventory category not found"
// @Router       /inventory-categories/{id}/toggle [patch]
// @Security     Bearer
func ToggleInventoryCategoryStatus(c *gin.Context) {
	id := c.Param("id")
	var category models.InventoryCategory

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory category not found"})
		return
	}

	category.IsActive = !category.IsActive

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    category,
	})
}
