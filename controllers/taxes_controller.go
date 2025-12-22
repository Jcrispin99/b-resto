package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTaxes godoc
// @Summary      Listar impuestos
// @Description  Obtiene lista de todos los impuestos configurados
// @Tags         taxes
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de taxes"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /taxes [get]
// @Security     Bearer
func GetTaxes(c *gin.Context) {
	var taxes []models.Tax

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if err := query.Find(&taxes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch taxes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": taxes})
}

// GetTax godoc
// @Summary      Obtener impuesto
// @Description  Obtiene un impuesto por ID
// @Tags         taxes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del impuesto"
// @Success      200  {object}  map[string]interface{}  "data: tax"
// @Failure      404  {object}  map[string]string       "error: Tax not found"
// @Router       /taxes/{id} [get]
// @Security     Bearer
func GetTax(c *gin.Context) {
	id := c.Param("id")
	var tax models.Tax

	if err := config.DB.First(&tax, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tax})
}

// CreateTax godoc
// @Summary      Crear impuesto
// @Description  Crea un nuevo impuesto
// @Tags         taxes
// @Accept       json
// @Produce      json
// @Param        tax  body  models.Tax  true  "Datos del impuesto"
// @Success      201  {object}  map[string]interface{}  "message y data creada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      409  {object}  map[string]string       "error: nombre duplicado"
// @Router       /taxes [post]
// @Security     Bearer
func CreateTax(c *gin.Context) {
	var tax models.Tax

	// Gin valida automáticamente según los binding tags del modelo
	if err := c.ShouldBindJSON(&tax); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Solo validaciones de negocio (unicidad)
	var existing models.Tax
	if err := config.DB.Where("name = ?", tax.Name).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Tax with this name already exists"})
		return
	}

	if err := config.DB.Create(&tax).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tax"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Tax created successfully",
		"data":    tax,
	})
}

// UpdateTax godoc
// @Summary      Actualizar impuesto
// @Description  Actualiza los datos de un impuesto existente
// @Tags         taxes
// @Accept       json
// @Produce      json
// @Param        id   path  int         true  "ID del impuesto"
// @Param        tax  body  models.Tax  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data actualizada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Tax not found"
// @Router       /taxes/{id} [put]
// @Security     Bearer
func UpdateTax(c *gin.Context) {
	id := c.Param("id")
	var tax models.Tax

	if err := config.DB.First(&tax, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax not found"})
		return
	}

	var updateData models.Tax
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&tax).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tax"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Tax updated successfully",
		"data":    tax,
	})
}

// DeleteTax godoc
// @Summary      Eliminar impuesto
// @Description  Elimina un impuesto (soft delete)
// @Tags         taxes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del impuesto"
// @Success      200  {object}  map[string]string  "message: Tax deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Tax not found"
// @Router       /taxes/{id} [delete]
// @Security     Bearer
func DeleteTax(c *gin.Context) {
	id := c.Param("id")
	var tax models.Tax

	if err := config.DB.First(&tax, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax not found"})
		return
	}

	if err := config.DB.Delete(&tax).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tax"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Tax deleted successfully"})
}

// ToggleTaxStatus godoc
// @Summary      Activar/Desactivar impuesto
// @Description  Cambia el estado is_active de un impuesto
// @Tags         taxes
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del impuesto"
// @Success      200  {object}  map[string]interface{}  "message y data con nuevo estado"
// @Failure      404  {object}  map[string]string       "error: Tax not found"
// @Router       /taxes/{id}/toggle [patch]
// @Security     Bearer
func ToggleTaxStatus(c *gin.Context) {
	id := c.Param("id")
	var tax models.Tax

	if err := config.DB.First(&tax, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax not found"})
		return
	}

	tax.IsActive = !tax.IsActive

	if err := config.DB.Save(&tax).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    tax,
	})
}
