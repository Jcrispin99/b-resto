package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetWarehouses godoc
// @Summary      Listar almacenes
// @Description  Obtiene lista de todos los almacenes
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        company_id query  int     false  "Filtrar por compañía/sucursal"
// @Success      200  {object}  map[string]interface{}  "data: array de warehouses"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /warehouses [get]
// @Security     Bearer
func GetWarehouses(c *gin.Context) {
	var warehouses []models.Warehouse

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if companyID := c.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if err := query.Preload("Company").Find(&warehouses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch warehouses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": warehouses})
}

// GetWarehouse godoc
// @Summary      Obtener almacén
// @Description  Obtiene un almacén por ID
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del almacén"
// @Success      200  {object}  map[string]interface{}  "data: warehouse"
// @Failure      404  {object}  map[string]string       "error: Warehouse not found"
// @Router       /warehouses/{id} [get]
// @Security     Bearer
func GetWarehouse(c *gin.Context) {
	id := c.Param("id")
	var warehouse models.Warehouse

	if err := config.DB.Preload("Company").First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": warehouse})
}

// CreateWarehouse godoc
// @Summary      Crear almacén
// @Description  Crea un nuevo almacén
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        warehouse  body  models.Warehouse  true  "Datos del almacén"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      409  {object}  map[string]string       "error: código duplicado"
// @Router       /warehouses [post]
// @Security     Bearer
func CreateWarehouse(c *gin.Context) {
	var warehouse models.Warehouse

	if err := c.ShouldBindJSON(&warehouse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar unicidad del código
	var existing models.Warehouse
	if err := config.DB.Where("code = ?", warehouse.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Warehouse with this code already exists"})
		return
	}

	if err := config.DB.Create(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create warehouse"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Warehouse created successfully",
		"data":    warehouse,
	})
}

// UpdateWarehouse godoc
// @Summary      Actualizar almacén
// @Description  Actualiza los datos de un almacén existente
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id         path  int               true  "ID del almacén"
// @Param        warehouse  body  models.Warehouse  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Warehouse not found"
// @Router       /warehouses/{id} [put]
// @Security     Bearer
func UpdateWarehouse(c *gin.Context) {
	id := c.Param("id")
	var warehouse models.Warehouse

	if err := config.DB.First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	var updateData models.Warehouse
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&warehouse).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update warehouse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Warehouse updated successfully",
		"data":    warehouse,
	})
}

// DeleteWarehouse godoc
// @Summary      Eliminar almacén
// @Description  Elimina un almacén (soft delete)
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del almacén"
// @Success      200  {object}  map[string]string  "message: Warehouse deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Warehouse not found"
// @Router       /warehouses/{id} [delete]
// @Security     Bearer
func DeleteWarehouse(c *gin.Context) {
	id := c.Param("id")
	var warehouse models.Warehouse

	if err := config.DB.First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	if err := config.DB.Delete(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete warehouse"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Warehouse deleted successfully"})
}

// ToggleWarehouseStatus godoc
// @Summary      Activar/Desactivar almacén
// @Description  Cambia el estado is_active de un almacén
// @Tags         warehouses
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del almacén"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Warehouse not found"
// @Router       /warehouses/{id}/toggle [patch]
// @Security     Bearer
func ToggleWarehouseStatus(c *gin.Context) {
	id := c.Param("id")
	var warehouse models.Warehouse

	if err := config.DB.First(&warehouse, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Warehouse not found"})
		return
	}

	warehouse.IsActive = !warehouse.IsActive

	if err := config.DB.Save(&warehouse).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    warehouse,
	})
}
