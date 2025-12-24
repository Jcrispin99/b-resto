package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTableAreas godoc
// @Summary      Listar áreas de mesas
// @Description  Obtiene lista de todas las áreas de mesas
// @Tags         table-areas
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        company_id query  int     false  "Filtrar por compañía/sucursal"
// @Success      200  {object}  map[string]interface{}  "data: array de table areas"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /table-areas [get]
// @Security     Bearer
func GetTableAreas(c *gin.Context) {
	var areas []models.TableArea

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if companyID := c.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if err := query.Preload("Company").Preload("Tables").Find(&areas).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch table areas"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": areas})
}

// GetTableArea godoc
// @Summary      Obtener área de mesas
// @Description  Obtiene un área de mesas por ID
// @Tags         table-areas
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del área"
// @Success      200  {object}  map[string]interface{}  "data: table area"
// @Failure      404  {object}  map[string]string       "error: Table area not found"
// @Router       /table-areas/{id} [get]
// @Security     Bearer
func GetTableArea(c *gin.Context) {
	id := c.Param("id")
	var area models.TableArea

	if err := config.DB.Preload("Company").Preload("Tables").First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table area not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": area})
}

// CreateTableArea godoc
// @Summary      Crear área de mesas
// @Description  Crea una nueva área de mesas
// @Tags         table-areas
// @Accept       json
// @Produce      json
// @Param        area  body  models.TableArea  true  "Datos del área"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /table-areas [post]
// @Security     Bearer
func CreateTableArea(c *gin.Context) {
	var area models.TableArea

	if err := c.ShouldBindJSON(&area); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&area).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table area"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Table area created successfully",
		"data":    area,
	})
}

// UpdateTableArea godoc
// @Summary      Actualizar área de mesas
// @Description  Actualiza los datos de un área de mesas existente
// @Tags         table-areas
// @Accept       json
// @Produce      json
// @Param        id    path  int                true  "ID del área"
// @Param        area  body  models.TableArea   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Table area not found"
// @Router       /table-areas/{id} [put]
// @Security     Bearer
func UpdateTableArea(c *gin.Context) {
	id := c.Param("id")
	var area models.TableArea

	if err := config.DB.First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table area not found"})
		return
	}

	var updateData models.TableArea
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&area).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update table area"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Table area updated successfully",
		"data":    area,
	})
}

// DeleteTableArea godoc
// @Summary      Eliminar área de mesas
// @Description  Elimina un área de mesas (soft delete)
// @Tags         table-areas
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del área"
// @Success      200  {object}  map[string]string  "message: Table area deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Table area not found"
// @Router       /table-areas/{id} [delete]
// @Security     Bearer
func DeleteTableArea(c *gin.Context) {
	id := c.Param("id")
	var area models.TableArea

	if err := config.DB.First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table area not found"})
		return
	}

	if err := config.DB.Delete(&area).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table area"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table area deleted successfully"})
}

// ToggleTableAreaStatus godoc
// @Summary      Activar/Desactivar área de mesas
// @Description  Cambia el estado is_active de un área de mesas
// @Tags         table-areas
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del área"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Table area not found"
// @Router       /table-areas/{id}/toggle [patch]
// @Security     Bearer
func ToggleTableAreaStatus(c *gin.Context) {
	id := c.Param("id")
	var area models.TableArea

	if err := config.DB.First(&area, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table area not found"})
		return
	}

	area.IsActive = !area.IsActive

	if err := config.DB.Save(&area).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    area,
	})
}
