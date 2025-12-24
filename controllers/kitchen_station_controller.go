package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetKitchenStations godoc
// @Summary      Listar estaciones de cocina
// @Description  Obtiene lista de todas las estaciones de cocina
// @Tags         kitchen-stations
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        company_id query  int     false  "Filtrar por compañía/sucursal"
// @Success      200  {object}  map[string]interface{}  "data: array de kitchen stations"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /kitchen-stations [get]
// @Security     Bearer
func GetKitchenStations(c *gin.Context) {
	var stations []models.KitchenStation

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if companyID := c.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if err := query.Preload("Company").Find(&stations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kitchen stations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": stations})
}

// GetKitchenStation godoc
// @Summary      Obtener estación de cocina
// @Description  Obtiene una estación de cocina por ID
// @Tags         kitchen-stations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la estación"
// @Success      200  {object}  map[string]interface{}  "data: kitchen station"
// @Failure      404  {object}  map[string]string       "error: Kitchen station not found"
// @Router       /kitchen-stations/{id} [get]
// @Security     Bearer
func GetKitchenStation(c *gin.Context) {
	id := c.Param("id")
	var station models.KitchenStation

	if err := config.DB.Preload("Company").First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen station not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": station})
}

// CreateKitchenStation godoc
// @Summary      Crear estación de cocina
// @Description  Crea una nueva estación de cocina
// @Tags         kitchen-stations
// @Accept       json
// @Produce      json
// @Param        station  body  models.KitchenStation  true  "Datos de la estación"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /kitchen-stations [post]
// @Security     Bearer
func CreateKitchenStation(c *gin.Context) {
	var station models.KitchenStation

	if err := c.ShouldBindJSON(&station); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&station).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create kitchen station"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Kitchen station created successfully",
		"data":    station,
	})
}

// UpdateKitchenStation godoc
// @Summary      Actualizar estación de cocina
// @Description  Actualiza los datos de una estación de cocina existente
// @Tags         kitchen-stations
// @Accept       json
// @Produce      json
// @Param        id       path  int                     true  "ID de la estación"
// @Param        station  body  models.KitchenStation   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Kitchen station not found"
// @Router       /kitchen-stations/{id} [put]
// @Security     Bearer
func UpdateKitchenStation(c *gin.Context) {
	id := c.Param("id")
	var station models.KitchenStation

	if err := config.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen station not found"})
		return
	}

	var updateData models.KitchenStation
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&station).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update kitchen station"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Kitchen station updated successfully",
		"data":    station,
	})
}

// DeleteKitchenStation godoc
// @Summary      Eliminar estación de cocina
// @Description  Elimina una estación de cocina (soft delete)
// @Tags         kitchen-stations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la estación"
// @Success      200  {object}  map[string]string  "message: Kitchen station deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Kitchen station not found"
// @Router       /kitchen-stations/{id} [delete]
// @Security     Bearer
func DeleteKitchenStation(c *gin.Context) {
	id := c.Param("id")
	var station models.KitchenStation

	if err := config.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen station not found"})
		return
	}

	if err := config.DB.Delete(&station).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete kitchen station"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Kitchen station deleted successfully"})
}

// ToggleKitchenStationStatus godoc
// @Summary      Activar/Desactivar estación de cocina
// @Description  Cambia el estado is_active de una estación de cocina
// @Tags         kitchen-stations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la estación"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Kitchen station not found"
// @Router       /kitchen-stations/{id}/toggle [patch]
// @Security     Bearer
func ToggleKitchenStationStatus(c *gin.Context) {
	id := c.Param("id")
	var station models.KitchenStation

	if err := config.DB.First(&station, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen station not found"})
		return
	}

	station.IsActive = !station.IsActive

	if err := config.DB.Save(&station).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    station,
	})
}
