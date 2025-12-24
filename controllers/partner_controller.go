package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPartners godoc
// @Summary      Listar partners
// @Description  Obtiene lista de todos los partners (proveedores/clientes)
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        is_active   query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        is_supplier query  string  false  "Filtrar por proveedor"      Enums(true, false)
// @Param        is_customer query  string  false  "Filtrar por cliente"        Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de partners"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /partners [get]
// @Security     Bearer
func GetPartners(c *gin.Context) {
	var partners []models.Partner

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if isSupplier := c.Query("is_supplier"); isSupplier != "" {
		query = query.Where("is_supplier = ?", isSupplier == "true")
	}
	if isCustomer := c.Query("is_customer"); isCustomer != "" {
		query = query.Where("is_customer = ?", isCustomer == "true")
	}

	if err := query.Find(&partners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch partners"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": partners})
}

// GetPartner godoc
// @Summary      Obtener partner
// @Description  Obtiene un partner por ID
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del partner"
// @Success      200  {object}  map[string]interface{}  "data: partner"
// @Failure      404  {object}  map[string]string       "error: Partner not found"
// @Router       /partners/{id} [get]
// @Security     Bearer
func GetPartner(c *gin.Context) {
	id := c.Param("id")
	var partner models.Partner

	if err := config.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": partner})
}

// CreatePartner godoc
// @Summary      Crear partner
// @Description  Crea un nuevo partner (proveedor/cliente)
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        partner  body  models.Partner  true  "Datos del partner"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /partners [post]
// @Security     Bearer
func CreatePartner(c *gin.Context) {
	var partner models.Partner

	if err := c.ShouldBindJSON(&partner); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create partner"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Partner created successfully",
		"data":    partner,
	})
}

// UpdatePartner godoc
// @Summary      Actualizar partner
// @Description  Actualiza los datos de un partner existente
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        id       path  int             true  "ID del partner"
// @Param        partner  body  models.Partner  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Partner not found"
// @Router       /partners/{id} [put]
// @Security     Bearer
func UpdatePartner(c *gin.Context) {
	id := c.Param("id")
	var partner models.Partner

	if err := config.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	var updateData models.Partner
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&partner).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update partner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Partner updated successfully",
		"data":    partner,
	})
}

// DeletePartner godoc
// @Summary      Eliminar partner
// @Description  Elimina un partner (soft delete)
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del partner"
// @Success      200  {object}  map[string]string  "message: Partner deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Partner not found"
// @Router       /partners/{id} [delete]
// @Security     Bearer
func DeletePartner(c *gin.Context) {
	id := c.Param("id")
	var partner models.Partner

	if err := config.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	if err := config.DB.Delete(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete partner"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Partner deleted successfully"})
}

// TogglePartnerStatus godoc
// @Summary      Activar/Desactivar partner
// @Description  Cambia el estado is_active de un partner
// @Tags         partners
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del partner"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Partner not found"
// @Router       /partners/{id}/toggle [patch]
// @Security     Bearer
func TogglePartnerStatus(c *gin.Context) {
	id := c.Param("id")
	var partner models.Partner

	if err := config.DB.First(&partner, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Partner not found"})
		return
	}

	partner.IsActive = !partner.IsActive

	if err := config.DB.Save(&partner).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    partner,
	})
}

// GetSuppliers godoc
// @Summary      Obtener proveedores
// @Description  Obtiene todos los partners marcados como proveedores
// @Tags         partners
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "data: array de suppliers"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /partners/suppliers [get]
// @Security     Bearer
func GetSuppliers(c *gin.Context) {
	var partners []models.Partner

	if err := config.DB.Where("is_supplier = ?", true).Find(&partners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch suppliers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": partners})
}

// GetCustomers godoc
// @Summary      Obtener clientes
// @Description  Obtiene todos los partners marcados como clientes
// @Tags         partners
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "data: array de customers"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /partners/customers [get]
// @Security     Bearer
func GetCustomers(c *gin.Context) {
	var partners []models.Partner

	if err := config.DB.Where("is_customer = ?", true).Find(&partners).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch customers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": partners})
}
