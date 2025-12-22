package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUnits obtiene todas las unidades
func GetUnits(c *gin.Context) {
	var units []models.Unit

	// // Filtro opcional por is_active
	// isActiveParam := c.Query("is_active")
	// query := config.DB

	// if isActiveParam != "" {
	// 	isActive, _ := strconv.ParseBool(isActiveParam)
	// 	query = query.Where("is_active = ?", isActive)
	// }

	// if err := query.Find(&units).Error; err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch units"})
	// 	return
	// }

	c.JSON(http.StatusOK, gin.H{
		"data": units,
	})
}

// GetUnit obtiene una unidad por ID
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

// CreateUnit crea una nueva unidad
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

// UpdateUnit actualiza una unidad existente
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

// DeleteUnit elimina una unidad (soft delete)
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

// ToggleUnitStatus activa/desactiva una unidad
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
