package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPOSTerminals godoc
// @Summary      Listar terminales POS
// @Description  Obtiene lista de todos los terminales POS
// @Tags         pos
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Param        company_id query  int     false  "Filtrar por compañía"
// @Success      200  {object}  map[string]interface{}  "data: array de pos terminals"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /pos [get]
// @Security     Bearer
func GetPOSTerminals(c *gin.Context) {
	var terminals []models.POS

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}
	if companyID := c.Query("company_id"); companyID != "" {
		query = query.Where("company_id = ?", companyID)
	}

	if err := query.Preload("Company").Find(&terminals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch POS terminals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": terminals})
}

// GetPOSTerminal godoc
// @Summary      Obtener terminal POS
// @Description  Obtiene un terminal POS por ID
// @Tags         pos
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del terminal"
// @Success      200  {object}  map[string]interface{}  "data: pos terminal"
// @Failure      404  {object}  map[string]string       "error: POS terminal not found"
// @Router       /pos/{id} [get]
// @Security     Bearer
func GetPOSTerminal(c *gin.Context) {
	id := c.Param("id")
	var terminal models.POS

	if err := config.DB.Preload("Company").First(&terminal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POS terminal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": terminal})
}

// CreatePOSTerminal godoc
// @Summary      Crear terminal POS
// @Description  Crea un nuevo terminal POS
// @Tags         pos
// @Accept       json
// @Produce      json
// @Param        pos  body  models.POS  true  "Datos del terminal"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /pos [post]
// @Security     Bearer
func CreatePOSTerminal(c *gin.Context) {
	var terminal models.POS

	if err := c.ShouldBindJSON(&terminal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&terminal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create POS terminal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "POS terminal created successfully",
		"data":    terminal,
	})
}

// UpdatePOSTerminal godoc
// @Summary      Actualizar terminal POS
// @Description  Actualiza los datos de un terminal POS existente
// @Tags         pos
// @Accept       json
// @Produce      json
// @Param        id   path  int         true  "ID del terminal"
// @Param        pos  body  models.POS  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: POS terminal not found"
// @Router       /pos/{id} [put]
// @Security     Bearer
func UpdatePOSTerminal(c *gin.Context) {
	id := c.Param("id")
	var terminal models.POS

	if err := config.DB.First(&terminal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POS terminal not found"})
		return
	}

	var updateData models.POS
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&terminal).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update POS terminal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "POS terminal updated successfully",
		"data":    terminal,
	})
}

// TogglePOSStatus godoc
// @Summary      Activar/Desactivar terminal POS
// @Description  Cambia el estado is_active de un terminal POS
// @Tags         pos
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del terminal"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: POS terminal not found"
// @Router       /pos/{id}/toggle [patch]
// @Security     Bearer
func TogglePOSStatus(c *gin.Context) {
	id := c.Param("id")
	var terminal models.POS

	if err := config.DB.First(&terminal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POS terminal not found"})
		return
	}

	terminal.IsActive = !terminal.IsActive

	if err := config.DB.Save(&terminal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    terminal,
	})
}
