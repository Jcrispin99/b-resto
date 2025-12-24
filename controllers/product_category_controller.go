package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProductCategories godoc
// @Summary      Listar categorías de productos
// @Description  Obtiene lista de todas las categorías de productos/menú con jerarquía
// @Tags         product-categories
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        parent_id  query  int     false  "Filtrar por categoría padre (null para raíz)"
// @Success      200  {object}  map[string]interface{}  "data: array de product categories"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /product-categories [get]
// @Security     Bearer
func GetProductCategories(c *gin.Context) {
	var categories []models.ProductCategory

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if parentID := c.Query("parent_id"); parentID != "" {
		if parentID == "null" {
			query = query.Where("parent_id IS NULL")
		} else {
			query = query.Where("parent_id = ?", parentID)
		}
	}

	// Cargar jerarquía completa
	query = query.Preload("Parent").Preload("Children")

	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product categories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GetProductCategory godoc
// @Summary      Obtener categoría de producto
// @Description  Obtiene una categoría de producto por ID con su jerarquía
// @Tags         product-categories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la categoría"
// @Success      200  {object}  map[string]interface{}  "data: product category"
// @Failure      404  {object}  map[string]string       "error: Product category not found"
// @Router       /product-categories/{id} [get]
// @Security     Bearer
func GetProductCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.ProductCategory

	if err := config.DB.Preload("Parent").Preload("Children").First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

// CreateProductCategory godoc
// @Summary      Crear categoría de producto
// @Description  Crea una nueva categoría de producto/menú
// @Tags         product-categories
// @Accept       json
// @Produce      json
// @Param        category  body  models.ProductCategory  true  "Datos de la categoría"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /product-categories [post]
// @Security     Bearer
func CreateProductCategory(c *gin.Context) {
	var category models.ProductCategory

	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar parent_id si se proporciona
	if category.ParentID != nil {
		var parent models.ProductCategory
		if err := config.DB.First(&parent, *category.ParentID).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Parent category not found"})
			return
		}
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product category"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product category created successfully",
		"data":    category,
	})
}

// UpdateProductCategory godoc
// @Summary      Actualizar categoría de producto
// @Description  Actualiza los datos de una categoría de producto existente
// @Tags         product-categories
// @Accept       json
// @Produce      json
// @Param        id        path  int                      true  "ID de la categoría"
// @Param        category  body  models.ProductCategory   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Product category not found"
// @Router       /product-categories/{id} [put]
// @Security     Bearer
func UpdateProductCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.ProductCategory

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
		return
	}

	var updateData models.ProductCategory
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product category updated successfully",
		"data":    category,
	})
}

// DeleteProductCategory godoc
// @Summary      Eliminar categoría de producto
// @Description  Elimina una categoría de producto (soft delete). No permite si tiene subcategorías.
// @Tags         product-categories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la categoría"
// @Success      200  {object}  map[string]string  "message: Product category deleted successfully"
// @Failure      400  {object}  map[string]string  "error: tiene subcategorías"
// @Failure      404  {object}  map[string]string  "error: Product category not found"
// @Router       /product-categories/{id} [delete]
// @Security     Bearer
func DeleteProductCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.ProductCategory

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
		return
	}

	// Verificar si tiene subcategorías
	var childCount int64
	config.DB.Model(&models.ProductCategory{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete category with subcategories"})
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product category"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product category deleted successfully"})
}

// ToggleProductCategoryStatus godoc
// @Summary      Activar/Desactivar categoría de producto
// @Description  Cambia el estado is_active de una categoría de producto
// @Tags         product-categories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la categoría"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Product category not found"
// @Router       /product-categories/{id}/toggle [patch]
// @Security     Bearer
func ToggleProductCategoryStatus(c *gin.Context) {
	id := c.Param("id")
	var category models.ProductCategory

	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product category not found"})
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
