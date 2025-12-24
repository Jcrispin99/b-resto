package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetPOSSessions godoc
// @Summary      Listar sesiones POS
// @Description  Obtiene lista de todas las sesiones de caja
// @Tags         pos-sessions
// @Accept       json
// @Produce      json
// @Param        pos_id  query  int     false  "Filtrar por terminal POS"
// @Param        state   query  string  false  "Filtrar por estado"  Enums(opened, closed)
// @Success      200  {object}  map[string]interface{}  "data: array de pos sessions"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /pos-sessions [get]
// @Security     Bearer
func GetPOSSessions(c *gin.Context) {
	var sessions []models.POSSession

	query := config.DB
	if posID := c.Query("pos_id"); posID != "" {
		query = query.Where("pos_id = ?", posID)
	}
	if state := c.Query("state"); state != "" {
		query = query.Where("state = ?", state)
	}

	if err := query.Preload("POS").Preload("User").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch POS sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// GetPOSSession godoc
// @Summary      Obtener sesión POS
// @Description  Obtiene una sesión de caja por ID
// @Tags         pos-sessions
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la sesión"
// @Success      200  {object}  map[string]interface{}  "data: pos session"
// @Failure      404  {object}  map[string]string       "error: POS session not found"
// @Router       /pos-sessions/{id} [get]
// @Security     Bearer
func GetPOSSession(c *gin.Context) {
	id := c.Param("id")
	var session models.POSSession

	if err := config.DB.Preload("POS").Preload("User").First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POS session not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": session})
}

// OpenPOSSession godoc
// @Summary      Abrir sesión de caja
// @Description  Abre una nueva sesión de caja en un terminal POS
// @Tags         pos-sessions
// @Accept       json
// @Produce      json
// @Param        session  body  models.POSSession  true  "Datos de la sesión (pos_id, user_id, opening_cash)"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación o sesión ya abierta"
// @Router       /pos-sessions/open [post]
// @Security     Bearer
func OpenPOSSession(c *gin.Context) {
	var session models.POSSession

	if err := c.ShouldBindJSON(&session); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verificar que no haya una sesión abierta en este POS
	var existingSession models.POSSession
	if err := config.DB.Where("pos_id = ? AND state = ?", session.POSID, "opened").First(&existingSession).Error; err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "There is already an open session on this POS terminal"})
		return
	}

	// Establecer valores iniciales
	now := time.Now()
	session.Status = "opened"
	session.OpenedAt = now
	diff := 0.0
	session.Difference = &diff

	if err := config.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open POS session"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "POS session opened successfully",
		"data":    session,
	})
}

// ClosePOSSession godoc
// @Summary      Cerrar sesión de caja
// @Description  Cierra una sesión de caja abierta
// @Tags         pos-sessions
// @Accept       json
// @Produce      json
// @Param        id       path  int                      true  "ID de la sesión"
// @Param        request  body  map[string]interface{}   true  "closing_cash: efectivo al cerrar"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación o sesión ya cerrada"
// @Failure      404  {object}  map[string]string       "error: POS session not found"
// @Router       /pos-sessions/{id}/close [patch]
// @Security     Bearer
func ClosePOSSession(c *gin.Context) {
	id := c.Param("id")
	var session models.POSSession

	if err := config.DB.First(&session, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "POS session not found"})
		return
	}

	if session.Status == "closed" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is already closed"})
		return
	}

	var request struct {
		ClosingCash float64 `json:"closing_cash" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	session.Status = "closed"
	session.ClosedAt = &now
	session.ClosingBalance = &request.ClosingCash

	// Calcular diferencia
	expectedCash := session.OpeningBalance
	diff := request.ClosingCash - expectedCash
	session.Difference = &diff

	if err := config.DB.Save(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close POS session"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "POS session closed successfully",
		"data":    session,
	})
}

// GetActivePOSSessions godoc
// @Summary      Obtener sesiones activas
// @Description  Obtiene todas las sesiones de caja abiertas
// @Tags         pos-sessions
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "data: array de sesiones abiertas"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /pos-sessions/active [get]
// @Security     Bearer
func GetActivePOSSessions(c *gin.Context) {
	var sessions []models.POSSession

	if err := config.DB.Where("state = ?", "opened").
		Preload("POS").Preload("User").
		Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch active sessions"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": sessions})
}

// GetSessionMovements godoc
// @Summary      Obtener movimientos de sesión
// @Description  Obtiene todos los movimientos de efectivo de una sesión
// @Tags         pos-sessions
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la sesión"
// @Success      200  {object}  map[string]interface{}  "data: array de cash movements"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /pos-sessions/{id}/movements [get]
// @Security     Bearer
func GetSessionMovements(c *gin.Context) {
	sessionID := c.Param("id")
	var movements []models.CashMovement

	if err := config.DB.Where("pos_session_id = ?", sessionID).
		Find(&movements).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movements"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": movements})
}
