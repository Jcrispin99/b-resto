package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetProductAttributeValues godoc
// @Summary      Listar valores de atributo
// @Description  Obtiene lista de todos los valores de atributos de producto
// @Tags         product-attribute-values
// @Accept       json
// @Produce      json
// @Param        attribute_id  query  int  false  "Filtrar por atributo"
// @Success      200  {object}  map[string]interface{}  "data: array de attribute values"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /product-attribute-values [get]
// @Security     Bearer
func GetProductAttributeValues(c *gin.Context) {
	var values []models.ProductAttributeValue

	query := config.DB
	if attributeID := c.Query("attribute_id"); attributeID != "" {
		query = query.Where("attribute_id = ?", attributeID)
	}

	if err := query.Preload("Attribute").Find(&values).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch attribute values"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": values})
}

// GetProductAttributeValue godoc
// @Summary      Obtener valor de atributo
// @Description  Obtiene un valor de atributo por ID
// @Tags         product-attribute-values
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del valor"
// @Success      200  {object}  map[string]interface{}  "data: attribute value"
// @Failure      404  {object}  map[string]string       "error: Attribute value not found"
// @Router       /product-attribute-values/{id} [get]
// @Security     Bearer
func GetProductAttributeValue(c *gin.Context) {
	id := c.Param("id")
	var value models.ProductAttributeValue

	if err := config.DB.Preload("Attribute").First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute value not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": value})
}

// CreateProductAttributeValue godoc
// @Summary      Crear valor de atributo
// @Description  Crea un nuevo valor de atributo
// @Tags         product-attribute-values
// @Accept       json
// @Produce      json
// @Param        value  body  models.ProductAttributeValue  true  "Datos del valor"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /product-attribute-values [post]
// @Security     Bearer
func CreateProductAttributeValue(c *gin.Context) {
	var value models.ProductAttributeValue

	if err := c.ShouldBindJSON(&value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create attribute value"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Attribute value created successfully",
		"data":    value,
	})
}

// UpdateProductAttributeValue godoc
// @Summary      Actualizar valor de atributo
// @Description  Actualiza los datos de un valor de atributo existente
// @Tags         product-attribute-values
// @Accept       json
// @Produce      json
// @Param        id     path  int                            true  "ID del valor"
// @Param        value  body  models.ProductAttributeValue   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Attribute value not found"
// @Router       /product-attribute-values/{id} [put]
// @Security     Bearer
func UpdateProductAttributeValue(c *gin.Context) {
	id := c.Param("id")
	var value models.ProductAttributeValue

	if err := config.DB.First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute value not found"})
		return
	}

	var updateData models.ProductAttributeValue
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&value).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update attribute value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Attribute value updated successfully",
		"data":    value,
	})
}

// DeleteProductAttributeValue godoc
// @Summary      Eliminar valor de atributo
// @Description  Elimina un valor de atributo (soft delete)
// @Tags         product-attribute-values
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del valor"
// @Success      200  {object}  map[string]string  "message: Attribute value deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Attribute value not found"
// @Router       /product-attribute-values/{id} [delete]
// @Security     Bearer
func DeleteProductAttributeValue(c *gin.Context) {
	id := c.Param("id")
	var value models.ProductAttributeValue

	if err := config.DB.First(&value, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Attribute value not found"})
		return
	}

	if err := config.DB.Delete(&value).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete attribute value"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Attribute value deleted successfully"})
}
