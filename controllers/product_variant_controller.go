package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProductVariants godoc
// @Summary      Listar variantes de producto
// @Description  Obtiene lista de todas las variantes/SKUs de productos
// @Tags         product-variants
// @Accept       json
// @Produce      json
// @Param        template_id  query  int     false  "Filtrar por template"
// @Param        is_active    query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de product variants"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /product-variants [get]
// @Security     Bearer
func GetProductVariants(c *gin.Context) {
	var variants []models.ProductProduct

	query := config.DB
	if templateID := c.Query("template_id"); templateID != "" {
		query = query.Where("template_id = ?", templateID)
	}
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if err := query.Preload("Template").Find(&variants).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product variants"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variants})
}

// GetProductVariant godoc
// @Summary      Obtener variante de producto
// @Description  Obtiene una variante/SKU por ID
// @Tags         product-variants
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la variante"
// @Success      200  {object}  map[string]interface{}  "data: product variant"
// @Failure      404  {object}  map[string]string       "error: Product variant not found"
// @Router       /product-variants/{id} [get]
// @Security     Bearer
func GetProductVariant(c *gin.Context) {
	id := c.Param("id")
	var variant models.ProductProduct

	if err := config.DB.Preload("Template").First(&variant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": variant})
}

// CreateProductVariant godoc
// @Summary      Crear variante de producto
// @Description  Crea una nueva variante/SKU de producto
// @Tags         product-variants
// @Accept       json
// @Produce      json
// @Param        variant  body  models.ProductProduct  true  "Datos de la variante"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /product-variants [post]
// @Security     Bearer
func CreateProductVariant(c *gin.Context) {
	var variant models.ProductProduct

	if err := c.ShouldBindJSON(&variant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product variant"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Product variant created successfully",
		"data":    variant,
	})
}

// UpdateProductVariant godoc
// @Summary      Actualizar variante de producto
// @Description  Actualiza los datos de una variante/SKU existente
// @Tags         product-variants
// @Accept       json
// @Produce      json
// @Param        id       path  int                     true  "ID de la variante"
// @Param        variant  body  models.ProductProduct   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Product variant not found"
// @Router       /product-variants/{id} [put]
// @Security     Bearer
func UpdateProductVariant(c *gin.Context) {
	id := c.Param("id")
	var variant models.ProductProduct

	if err := config.DB.First(&variant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
		return
	}

	var updateData models.ProductProduct
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&variant).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Product variant updated successfully",
		"data":    variant,
	})
}

// DeleteProductVariant godoc
// @Summary      Eliminar variante de producto
// @Description  Elimina una variante/SKU (soft delete)
// @Tags         product-variants
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la variante"
// @Success      200  {object}  map[string]string  "message: Product variant deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Product variant not found"
// @Router       /product-variants/{id} [delete]
// @Security     Bearer
func DeleteProductVariant(c *gin.Context) {
	id := c.Param("id")
	var variant models.ProductProduct

	if err := config.DB.First(&variant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
		return
	}

	if err := config.DB.Delete(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product variant"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product variant deleted successfully"})
}

// ToggleProductVariantStatus godoc
// @Summary      Activar/Desactivar variante
// @Description  Cambia el estado is_active de una variante de producto
// @Tags         product-variants
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la variante"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Product variant not found"
// @Router       /product-variants/{id}/toggle [patch]
// @Security     Bearer
func ToggleProductVariantStatus(c *gin.Context) {
	id := c.Param("id")
	var variant models.ProductProduct

	if err := config.DB.First(&variant, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product variant not found"})
		return
	}

	variant.IsActive = !variant.IsActive

	if err := config.DB.Save(&variant).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    variant,
	})
}
