package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTables godoc
// @Summary      Listar mesas
// @Description  Obtiene lista de todas las mesas
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        area_id    query  int     false  "Filtrar por área"
// @Success      200  {object}  map[string]interface{}  "data: array de tables"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /tables [get]
// @Security     Bearer
func GetTables(c *gin.Context) {
	var tables []models.Table

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if areaID := c.Query("area_id"); areaID != "" {
		query = query.Where("area_id = ?", areaID)
	}

	if err := query.Preload("Area").Preload("Company").Find(&tables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tables})
}

// GetTable godoc
// @Summary      Obtener mesa
// @Description  Obtiene una mesa por ID
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la mesa"
// @Success      200  {object}  map[string]interface{}  "data: table"
// @Failure      404  {object}  map[string]string       "error: Table not found"
// @Router       /tables/{id} [get]
// @Security     Bearer
func GetTable(c *gin.Context) {
	id := c.Param("id")
	var table models.Table

	if err := config.DB.Preload("Area").Preload("Company").First(&table, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": table})
}

// CreateTable godoc
// @Summary      Crear mesa
// @Description  Crea una nueva mesa
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        table  body  models.Table  true  "Datos de la mesa"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /tables [post]
// @Security     Bearer
func CreateTable(c *gin.Context) {
	var table models.Table

	if err := c.ShouldBindJSON(&table); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&table).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create table"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Table created successfully",
		"data":    table,
	})
}

// UpdateTable godoc
// @Summary      Actualizar mesa
// @Description  Actualiza los datos de una mesa existente
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id     path  int           true  "ID de la mesa"
// @Param        table  body  models.Table  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Table not found"
// @Router       /tables/{id} [put]
// @Security     Bearer
func UpdateTable(c *gin.Context) {
	id := c.Param("id")
	var table models.Table

	if err := config.DB.First(&table, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	var updateData models.Table
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&table).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Table updated successfully",
		"data":    table,
	})
}

// DeleteTable godoc
// @Summary      Eliminar mesa
// @Description  Elimina una mesa (soft delete)
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la mesa"
// @Success      200  {object}  map[string]string  "message: Table deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Table not found"
// @Router       /tables/{id} [delete]
// @Security     Bearer
func DeleteTable(c *gin.Context) {
	id := c.Param("id")
	var table models.Table

	if err := config.DB.First(&table, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	if err := config.DB.Delete(&table).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete table"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Table deleted successfully"})
}

// ToggleTableStatus godoc
// @Summary      Activar/Desactivar mesa
// @Description  Cambia el estado is_active de una mesa
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la mesa"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Table not found"
// @Router       /tables/{id}/toggle [patch]
// @Security     Bearer
func ToggleTableStatus(c *gin.Context) {
	id := c.Param("id")
	var table models.Table

	if err := config.DB.First(&table, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Table not found"})
		return
	}

	table.IsActive = !table.IsActive

	if err := config.DB.Save(&table).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    table,
	})
}

// GetTablesByArea godoc
// @Summary      Obtener mesas por área
// @Description  Obtiene todas las mesas de un área específica
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        area_id  path  int  true  "ID del área"
// @Success      200  {object}  map[string]interface{}  "data: array de tables"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /tables/area/{area_id} [get]
// @Security     Bearer
func GetTablesByArea(c *gin.Context) {
	areaID := c.Param("area_id")
	var tables []models.Table

	if err := config.DB.Where("area_id = ?", areaID).Preload("Area").Find(&tables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tables})
}

// GetAvailableTables godoc
// @Summary      Obtener mesas disponibles
// @Description  Obtiene todas las mesas activas sin órdenes pendientes
// @Tags         tables
// @Accept       json
// @Produce      json
// @Param        area_id  query  int  false  "Filtrar por área"
// @Success      200  {object}  map[string]interface{}  "data: array de tables disponibles"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /tables/available [get]
// @Security     Bearer
func GetAvailableTables(c *gin.Context) {
	query := config.DB.Where("is_active = ?", true)

	if areaID := c.Query("area_id"); areaID != "" {
		query = query.Where("area_id = ?", areaID)
	}

	// Obtener mesas que NO tienen órdenes en estado 'pending' o 'confirmed'
	var tables []models.Table
	if err := query.Preload("Area").
		Where("id NOT IN (?)",
			config.DB.Table("orders").
				Select("table_id").
				Where("table_id IS NOT NULL").
				Where("status IN (?)", []string{"pending", "confirmed"}),
		).Find(&tables).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch available tables"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tables})
}
