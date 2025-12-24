package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProductAttributes godoc
// @Summary      Listar atributos de producto
// @Description  Obtiene lista de todos los atributos (Tamaño, Color, etc.)
// @Tags         product-attributes
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "data: array de product attributes"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /product-attributes [get]
// @Security     Bearer
func GetProductAttributes(c *gin.Context) {
	var attributes []models.ProductAttribute
	if err := config.DB.Preload("Values").Find(&attributes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product attributes"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": attributes})
}

// GetProductAttribute godoc
// @Summary      Obtener atributo de producto
// @Description  Obtiene un atributo de producto por ID con sus valores
// @Tags         product-attributes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del atributo"
// @Success      200  {object}  map[string]interface{}  "data: product attribute"
// @Failure      404  {object}  map[string]string       "error: Product attribute not found"
// @Router       /product-attributes/{id} [get]
// @Security     Bearer
func GetProductAttribute(c *gin.Context) {
	id := c.Param("id")
	var attribute models.ProductAttribute
	if err := config.DB.Preload("Values").First(&attribute, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product attribute not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": attribute})
}

// CreateProductAttribute godoc
// @Summary      Crear atributo de producto
// @Description  Crea un nuevo atributo de producto
// @Tags         product-attributes
// @Accept       json
// @Produce      json
// @Param        attribute  body  models.ProductAttribute  true  "Datos del atributo"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /product-attributes [post]
// @Security     Bearer
func CreateProductAttribute(c *gin.Context) {
	var attribute models.ProductAttribute
	if err := c.ShouldBindJSON(&attribute); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product attribute"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Product attribute created successfully",
		"data":    attribute,
	})
}

// UpdateProductAttribute godoc
// @Summary      Actualizar atributo de producto
// @Description  Actualiza los datos de un atributo de producto existente
// @Tags         product-attributes
// @Accept       json
// @Produce      json
// @Param        id         path  int                       true  "ID del atributo"
// @Param        attribute  body  models.ProductAttribute   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Product attribute not found"
// @Router       /product-attributes/{id} [put]
// @Security     Bearer
func UpdateProductAttribute(c *gin.Context) {
	id := c.Param("id")
	var attribute models.ProductAttribute
	if err := config.DB.First(&attribute, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product attribute not found"})
		return
	}
	var updateData models.ProductAttribute
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&attribute).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product attribute"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Product attribute updated successfully",
		"data":    attribute,
	})
}

// DeleteProductAttribute godoc
// @Summary      Eliminar atributo de producto
// @Description  Elimina un atributo de producto (soft delete)
// @Tags         product-attributes
// @Accept       json
// @Produce      json
// @Param         id  path  int  true  "ID del atributo"
// @Success      200  {object}  map[string]string  "message: Product attribute deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Product attribute not found"
// @Router       /product-attributes/{id} [delete]
// @Security     Bearer
func DeleteProductAttribute(c *gin.Context) {
	id := c.Param("id")
	var attribute models.ProductAttribute
	if err := config.DB.First(&attribute, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product attribute not found"})
		return
	}
	if err := config.DB.Delete(&attribute).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete product attribute"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product attribute deleted successfully"})
}
