package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSequences godoc
// @Summary      Listar secuencias
// @Description  Obtiene lista de todas las secuencias de numeración
// @Tags         sequences
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "data: array de sequences"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /sequences [get]
// @Security     Bearer
func GetSequences(c *gin.Context) {
	var sequences []models.Sequence
	if err := config.DB.Find(&sequences).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sequences"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sequences})
}

// GetSequence godoc
// @Summary      Obtener secuencia
// @Description  Obtiene una secuencia por ID
// @Tags         sequences
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la secuencia"
// @Success      200  {object}  map[string]interface{}  "data: sequence"
// @Failure      404  {object}  map[string]string       "error: Sequence not found"
// @Router       /sequences/{id} [get]
// @Security     Bearer
func GetSequence(c *gin.Context) {
	id := c.Param("id")
	var sequence models.Sequence
	if err := config.DB.First(&sequence, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sequence not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": sequence})
}

// CreateSequence godoc
// @Summary      Crear secuencia
// @Description  Crea una nueva secuencia de numeración
// @Tags         sequences
// @Accept       json
// @Produce      json
// @Param        sequence  body  models.Sequence  true  "Datos de la secuencia"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /sequences [post]
// @Security     Bearer
func CreateSequence(c *gin.Context) {
	var sequence models.Sequence
	if err := c.ShouldBindJSON(&sequence); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&sequence).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create sequence"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"message": "Sequence created successfully",
		"data":    sequence,
	})
}

// UpdateSequence godoc
// @Summary      Actualizar secuencia
// @Description  Actualiza los datos de una secuencia existente
// @Tags         sequences
// @Accept       json
// @Produce      json
// @Param        id        path  int              true  "ID de la secuencia"
// @Param        sequence  body  models.Sequence  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Sequence not found"
// @Router       /sequences/{id} [put]
// @Security     Bearer
func UpdateSequence(c *gin.Context) {
	id := c.Param("id")
	var sequence models.Sequence
	if err := config.DB.First(&sequence, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sequence not found"})
		return
	}
	var updateData models.Sequence
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Model(&sequence).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sequence"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Sequence updated successfully",
		"data":    sequence,
	})
}

// GetNextNumber godoc
// @Summary      Obtener siguiente número
// @Description  Obtiene el siguiente número de una secuencia por ID e incrementa el contador
// @Tags         sequences
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la secuencia"
// @Success      200  {object}  map[string]interface{}  "next_number: siguiente número formateado"
// @Failure      404  {object}  map[string]string       "error: Sequence not found"
// @Router       /sequences/{id}/next [get]
// @Security     Bearer
func GetNextNumber(c *gin.Context) {
	id := c.Param("id")
	var sequence models.Sequence

	// Buscar por ID con Journal
	if err := config.DB.Preload("Journal").First(&sequence, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Sequence not found"})
		return
	}

	// Obtener número actual
	currentNumber := sequence.NextNumber

	// Incrementar para el próximo
	sequence.NextNumber++

	// Formatear con padding
	paddedNumber := strconv.Itoa(currentNumber)
	if sequence.Padding > 0 {
		format := "%0" + strconv.Itoa(sequence.Padding) + "d"
		paddedNumber = fmt.Sprintf(format, currentNumber)
	}

	// Guardar el nuevo número
	if err := config.DB.Save(&sequence).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update sequence"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"next_number":   currentNumber,
		"formatted":     paddedNumber,
		"sequence_name": sequence.Name,
	})
}
