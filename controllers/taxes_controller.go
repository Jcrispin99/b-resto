package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTaxes obtiene todos los impuestos
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

// GetTax obtiene un impuesto por ID
func GetTax(c *gin.Context) {
	id := c.Param("id")
	var tax models.Tax

	if err := config.DB.First(&tax, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tax not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tax})
}

// CreateTax crea un nuevo impuesto
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

// UpdateTax actualiza un impuesto existente
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

// DeleteTax elimina un impuesto (soft delete)
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

// ToggleTaxStatus activa/desactiva un impuesto
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
