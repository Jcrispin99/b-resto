package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCashMovements godoc
// @Summary      Listar movimientos de efectivo
// @Description  Obtiene lista de movimientos de efectivo por sesión
// @Tags         cash-movements
// @Accept       json
// @Produce      json
// @Param        session_id  path  int  true  "ID de la sesión"
// @Success      200  {object}  map[string]interface{}  "data: array de cash movements"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /pos-sessions/{session_id}/cash-movements [get]
// @Security     Bearer
func GetCashMovements(c *gin.Context) {
	sessionID := c.Param("id")
	var movements []models.CashMovement

	if err := config.DB.Where("pos_session_id = ?", sessionID).
		Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cash movements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movements})
}

// CreateCashMovement godoc
// @Summary      Crear movimiento de efectivo
// @Description  Registra un nuevo movimiento de efectivo (ingreso/egreso)
// @Tags         cash-movements
// @Accept       json
// @Produce      json
// @Param        session_id  path  int                   true  "ID de la sesión"
// @Param        movement    body  models.CashMovement   true  "Datos del movimiento"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /pos-sessions/{session_id}/cash-movements [post]
// @Security     Bearer
func CreateCashMovement(c *gin.Context) {
	sessionID := c.Param("id")
	var movement models.CashMovement

	if err := c.ShouldBindJSON(&movement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar el SessionID del path
	var sessionIDUint uint
	if _, err := fmt.Sscanf(sessionID, "%d", &sessionIDUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session ID"})
		return
	}
	movement.POSSessionID = sessionIDUint

	if err := config.DB.Create(&movement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cash movement"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Cash movement created successfully",
		"data":    movement,
	})
}

// DeleteCashMovement godoc
// @Summary      Eliminar movimiento de efectivo
// @Description  Elimina un movimiento de efectivo de una sesión de caja
// @Tags         cash-movements
// @Accept       json
// @Produce      json
// @Param        session_id   path  int  true  "ID de la sesión"
// @Param        movement_id  path  int  true  "ID del movimiento"
// @Success      200  {object}  map[string]string  "message: Cash movement deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Cash movement not found"
// @Router       /pos-sessions/{session_id}/cash-movements/{movement_id} [delete]
// @Security     Bearer
func DeleteCashMovement(c *gin.Context) {
	movementID := c.Param("movement_id")
	var movement models.CashMovement

	if err := config.DB.First(&movement, movementID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cash movement not found"})
		return
	}

	if err := config.DB.Delete(&movement).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete cash movement"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cash movement deleted successfully"})
}
