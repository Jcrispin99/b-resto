package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetJournals godoc
// @Summary      Listar journals
// @Description  Obtiene lista de todos los diarios contables
// @Tags         journals
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de journals"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /journals [get]
// @Security     Bearer
func GetJournals(c *gin.Context) {
	var journals []models.Journal

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if err := query.Find(&journals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch journals"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": journals})
}

// GetJournal godoc
// @Summary      Obtener journal
// @Description  Obtiene un diario contable por ID
// @Tags         journals
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del journal"
// @Success      200  {object}  map[string]interface{}  "data: journal"
// @Failure      404  {object}  map[string]string       "error: Journal not found"
// @Router       /journals/{id} [get]
// @Security     Bearer
func GetJournal(c *gin.Context) {
	id := c.Param("id")
	var journal models.Journal

	if err := config.DB.First(&journal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Journal not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": journal})
}

// CreateJournal godoc
// @Summary      Crear journal
// @Description  Crea un nuevo diario contable
// @Tags         journals
// @Accept       json
// @Produce      json
// @Param        journal  body  models.Journal  true  "Datos del journal"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /journals [post]
// @Security     Bearer
func CreateJournal(c *gin.Context) {
	var journal models.Journal

	if err := c.ShouldBindJSON(&journal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&journal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create journal"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Journal created successfully",
		"data":    journal,
	})
}

// UpdateJournal godoc
// @Summary      Actualizar journal
// @Description  Actualiza los datos de un diario contable existente
// @Tags         journals
// @Accept       json
// @Produce      json
// @Param        id       path  int             true  "ID del journal"
// @Param        journal  body  models.Journal  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Journal not found"
// @Router       /journals/{id} [put]
// @Security     Bearer
func UpdateJournal(c *gin.Context) {
	id := c.Param("id")
	var journal models.Journal

	if err := config.DB.First(&journal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Journal not found"})
		return
	}

	var updateData models.Journal
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&journal).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update journal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Journal updated successfully",
		"data":    journal,
	})
}

// ToggleJournalStatus godoc
// @Summary      Activar/Desactivar journal
// @Description  Cambia el estado is_active de un diario contable
// @Tags         journals
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del journal"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Journal not found"
// @Router       /journals/{id}/toggle [patch]
// @Security     Bearer
func ToggleJournalStatus(c *gin.Context) {
	id := c.Param("id")
	var journal models.Journal

	if err := config.DB.First(&journal, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Journal not found"})
		return
	}

	journal.IsActive = !journal.IsActive

	if err := config.DB.Save(&journal).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    journal,
	})
}
