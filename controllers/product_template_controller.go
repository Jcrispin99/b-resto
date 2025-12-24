package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"b-resto/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProductTemplates godoc
// @Summary      Listar plantillas de producto
// @Description  Obtiene lista de todas las plantillas de producto (maestro)
// @Tags         product-templates
// @Accept       json
// @Produce      json
// @Param        is_active    query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        category_id  query  int     false  "Filtrar por categoría"
// @Param        can_be_sold  query  string  false  "Filtrar por vendible"       Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de product templates"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /product-templates [get]
// @Security     Bearer
func GetProductTemplates(c *gin.Context) {
	var templates []models.ProductTemplate

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if canBeSold := c.Query("can_be_sold"); canBeSold != "" {
		query = query.Where("can_be_sold = ?", canBeSold == "true")
	}

	if err := query.
		Preload("Category").
		Preload("Unit").
		Preload("KitchenStation").
		Preload("Variants").
		Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product templates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": templates})
}

// GetProductTemplate godoc
// @Summary      Obtener plantilla de producto
// @Description  Obtiene una plantilla de producto por ID con todas sus relaciones
// @Tags         product-templates
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la plantilla"
// @Success      200  {object}  map[string]interface{}  "data: product template"
// @Failure      404  {object}  map[string]string       "error: Product template not found"
// @Router       /product-templates/{id} [get]
// @Security     Bearer
func GetProductTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ProductTemplate

	if err := config.DB.
		Preload("Category").
		Preload("InventoryCategory").
		Preload("Unit").
		Preload("KitchenStation").
		Preload("Variants").
		First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product template not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}

// CreateProductTemplate godoc
// @Summary      Crear plantilla de producto
// @Description  Crea una nueva plantilla de producto con variante default automática
// @Tags         product-templates
// @Accept       json
// @Produce      json
// @Param        template  body  models.ProductTemplate  true  "Datos de la plantilla"
// @Success      201  {object}  map[string]interface{}  "message, template y default_variant"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /product-templates [post]
// @Security     Bearer
func CreateProductTemplate(c *gin.Context) {
	var template models.ProductTemplate

	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Usar servicio para crear template + variante default
	productService := services.NewProductService()
	createdTemplate, defaultVariant, err := productService.CreateTemplateWithDefaultVariant(&template)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":         "Product template created successfully with default variant",
		"template":        createdTemplate,
		"default_variant": defaultVariant,
	})
}

// UpdateProductTemplate godoc
// @Summary      Actualizar plantilla de producto
// @Description  Actualiza los datos de una plantilla de producto existente
// @Tags         product-templates
// @Accept       json
// @Produce      json
// @Param        id        path  int                      true  "ID de la plantilla"
// @Param        template  body  models.ProductTemplate   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Product template not found"
// @Router       /product-templates/{id} [put]
// @Security     Bearer
func UpdateProductTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ProductTemplate

	if err := config.DB.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product template not found"})
		return
	}

	var updateData models.ProductTemplate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&template).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product template updated successfully",
		"data":    template,
	})
}

// DeleteProductTemplate godoc
// @Summary      Eliminar plantilla de producto
// @Description  Elimina una plantilla de producto (soft delete)
// @Tags         product-templates
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la plantilla"
// @Success      200  {object}  map[string]string  "message: Product template deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Product template not found"
// @Router       /product-templates/{id} [delete]
// @Security     Bearer
func DeleteProductTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ProductTemplate

	if err := config.DB.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product template not found"})
		return
	}

	if err := config.DB.Delete(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product template"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product template deleted successfully"})
}

// ToggleProductTemplateStatus godoc
// @Summary      Activar/Desactivar plantilla
// @Description  Cambia el estado is_active de una plantilla de producto
// @Tags         product-templates
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la plantilla"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Product template not found"
// @Router       /product-templates/{id}/toggle [patch]
// @Security     Bearer
func ToggleProductTemplateStatus(c *gin.Context) {
	id := c.Param("id")
	var template models.ProductTemplate

	if err := config.DB.First(&template, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product template not found"})
		return
	}

	template.IsActive = !template.IsActive

	if err := config.DB.Save(&template).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    template,
	})
}
