package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCombos godoc
// @Summary      Listar combos
// @Description  Obtiene lista de todos los combos/paquetes de productos
// @Tags         combos
// @Accept       json
// @Produce      json
// @Param        is_active    query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        category_id  query  int     false  "Filtrar por categoría"
// @Success      200  {object}  map[string]interface{}  "data: array de combos"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /combos [get]
// @Security     Bearer
func GetCombos(c *gin.Context) {
	var combos []models.Combo

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if categoryID := c.Query("category_id"); categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	if err := query.Preload("Category").Preload("Items").Find(&combos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch combos"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": combos})
}

// GetCombo godoc
// @Summary      Obtener combo
// @Description  Obtiene un combo por ID con todos sus items
// @Tags         combos
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del combo"
// @Success      200  {object}  map[string]interface{}  "data: combo con items"
// @Failure      404  {object}  map[string]string       "error: Combo not found"
// @Router       /combos/{id} [get]
// @Security     Bearer
func GetCombo(c *gin.Context) {
	id := c.Param("id")
	var combo models.Combo

	if err := config.DB.Preload("Category").Preload("Items").First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": combo})
}

// CreateCombo godoc
// @Summary      Crear combo
// @Description  Crea un nuevo combo/paquete de productos
// @Tags         combos
// @Accept       json
// @Produce      json
// @Param        combo  body  models.Combo  true  "Datos del combo"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /combos [post]
// @Security     Bearer
func CreateCombo(c *gin.Context) {
	var combo models.Combo

	if err := c.ShouldBindJSON(&combo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&combo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create combo"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Combo created successfully",
		"data":    combo,
	})
}

// UpdateCombo godoc
// @Summary      Actualizar combo
// @Description  Actualiza los datos de un combo existente
// @Tags         combos
// @Accept       json
// @Produce      json
// @Param        id     path  int           true  "ID del combo"
// @Param        combo  body  models.Combo  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Combo not found"
// @Router       /combos/{id} [put]
// @Security     Bearer
func UpdateCombo(c *gin.Context) {
	id := c.Param("id")
	var combo models.Combo

	if err := config.DB.First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	var updateData models.Combo
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&combo).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update combo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Combo updated successfully",
		"data":    combo,
	})
}

// DeleteCombo godoc
// @Summary      Eliminar combo
// @Description  Elimina un combo (soft delete)
// @Tags         combos
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del combo"
// @Success      200  {object}  map[string]string  "message: Combo deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Combo not found"
// @Router       /combos/{id} [delete]
// @Security     Bearer
func DeleteCombo(c *gin.Context) {
	id := c.Param("id")
	var combo models.Combo

	if err := config.DB.First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	if err := config.DB.Delete(&combo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete combo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Combo deleted successfully"})
}

// ToggleComboStatus godoc
// @Summary      Activar/Desactivar combo
// @Description  Cambia el estado is_active de un combo
// @Tags         combos
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del combo"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Combo not found"
// @Router       /combos/{id}/toggle [patch]
// @Security     Bearer
func ToggleComboStatus(c *gin.Context) {
	id := c.Param("id")
	var combo models.Combo

	if err := config.DB.First(&combo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Combo not found"})
		return
	}

	combo.IsActive = !combo.IsActive

	if err := config.DB.Save(&combo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    combo,
	})
}
