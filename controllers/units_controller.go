package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUnits godoc
// @Summary      Listar unidades
// @Description  Obtiene lista de todas las unidades de medida
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de units"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /units [get]
// @Security     Bearer
func GetUnits(c *gin.Context) {
	var units []models.Unit

	// Filtro opcional por is_active
	isActiveParam := c.Query("is_active")
	query := config.DB

	if isActiveParam != "" {
		isActive, _ := strconv.ParseBool(isActiveParam)
		query = query.Where("is_active = ?", isActive)
	}

	if err := query.Find(&units).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch units"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": units,
	})
}

// GetUnit godoc
// @Summary      Obtener unidad
// @Description  Obtiene una unidad de medida por ID
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la unidad"
// @Success      200  {object}  map[string]interface{}  "data: unit"
// @Failure      404  {object}  map[string]string       "error: Unit not found"
// @Router       /units/{id} [get]
// @Security     Bearer
func GetUnit(c *gin.Context) {
	id := c.Param("id")
	var unit models.Unit

	if err := config.DB.First(&unit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": unit,
	})
}

// CreateUnit godoc
// @Summary      Crear unidad
// @Description  Crea una nueva unidad de medida
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        unit  body  models.Unit  true  "Datos de la unidad (name, abbreviation, type)"
// @Success      201  {object}  map[string]interface{}  "message y data creada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      409  {object}  map[string]string       "error: nombre o abreviación duplicada"
// @Router       /units [post]
// @Security     Bearer
func CreateUnit(c *gin.Context) {
	var unit models.Unit

	// Gin valida automáticamente según los binding tags del modelo
	if err := c.ShouldBindJSON(&unit); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Solo validaciones de negocio (unicidad)
	var existing models.Unit
	if err := config.DB.Where("name = ? OR abbreviation = ?", unit.Name, unit.Abbreviation).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Unit with this name or abbreviation already exists"})
		return
	}

	if err := config.DB.Create(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create unit"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Unit created successfully",
		"data":    unit,
	})
}

// UpdateUnit godoc
// @Summary      Actualizar unidad
// @Description  Actualiza los datos de una unidad existente
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        id    path  int          true  "ID de la unidad"
// @Param        unit  body  models.Unit  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data actualizada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Unit not found"
// @Router       /units/{id} [put]
// @Security     Bearer
func UpdateUnit(c *gin.Context) {
	id := c.Param("id")
	var unit models.Unit

	// Buscar la unidad
	if err := config.DB.First(&unit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	// Parsear los datos de actualización
	var updateData models.Unit
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Actualizar solo los campos proporcionados
	if err := config.DB.Model(&unit).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update unit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unit updated successfully",
		"data":    unit,
	})
}

// DeleteUnit godoc
// @Summary      Eliminar unidad
// @Description  Elimina una unidad de medida (soft delete)
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la unidad"
// @Success      200  {object}  map[string]string  "message: Unit deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Unit not found"
// @Router       /units/{id} [delete]
// @Security     Bearer
func DeleteUnit(c *gin.Context) {
	id := c.Param("id")
	var unit models.Unit

	if err := config.DB.First(&unit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	// Soft delete (GORM automáticamente usa DeletedAt)
	if err := config.DB.Delete(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete unit"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unit deleted successfully",
	})
}

// ToggleUnitStatus godoc
// @Summary      Activar/Desactivar unidad
// @Description  Cambia el estado is_active de una unidad
// @Tags         units
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la unidad"
// @Success      200  {object}  map[string]interface{}  "message y data con nuevo estado"
// @Failure      404  {object}  map[string]string       "error: Unit not found"
// @Router       /units/{id}/toggle [patch]
// @Security     Bearer
func ToggleUnitStatus(c *gin.Context) {
	id := c.Param("id")
	var unit models.Unit

	if err := config.DB.First(&unit, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unit not found"})
		return
	}

	// Toggle el estado
	unit.IsActive = !unit.IsActive

	if err := config.DB.Save(&unit).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Unit status updated successfully",
		"data":    unit,
	})
}
