package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// GetReservations godoc
// @Summary      Listar reservas
// @Description  Obtiene lista de todas las reservas
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        status  query  string  false  "Filtrar por estado"  Enums(pending, confirmed, seated, cancelled, no_show)
// @Param        date    query  string  false  "Filtrar por fecha (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}  "data: array de reservations"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /reservations [get]
// @Security     Bearer
func GetReservations(c *gin.Context) {
	var reservations []models.Reservation

	query := config.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if date := c.Query("date"); date != "" {
		query = query.Where("reservation_date = ?", date)
	}

	if err := query.Preload("Company").Preload("Partner").Preload("Table").Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reservations})
}

// GetReservation godoc
// @Summary      Obtener reserva
// @Description  Obtiene una reserva por ID
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la reserva"
// @Success      200  {object}  map[string]interface{}  "data: reservation"
// @Failure      404  {object}  map[string]string       "error: Reservation not found"
// @Router       /reservations/{id} [get]
// @Security     Bearer
func GetReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation models.Reservation

	if err := config.DB.Preload("Company").Preload("Partner").Preload("Table").First(&reservation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reservation})
}

// CreateReservation godoc
// @Summary      Crear reserva
// @Description  Crea una nueva reserva
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        reservation  body  models.Reservation  true  "Datos de la reserva"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /reservations [post]
// @Security     Bearer
func CreateReservation(c *gin.Context) {
	var reservation models.Reservation

	if err := c.ShouldBindJSON(&reservation); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Estado inicial
	if reservation.Status == "" {
		reservation.Status = "pending"
	}

	if err := config.DB.Create(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Reservation created successfully",
		"data":    reservation,
	})
}

// UpdateReservation godoc
// @Summary      Actualizar reserva
// @Description  Actualiza los datos de una reserva existente
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id           path  int                  true  "ID de la reserva"
// @Param        reservation  body  models.Reservation   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Reservation not found"
// @Router       /reservations/{id} [put]
// @Security     Bearer
func UpdateReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation models.Reservation

	if err := config.DB.First(&reservation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	var updateData models.Reservation
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&reservation).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation updated successfully",
		"data":    reservation,
	})
}

// ConfirmReservation godoc
// @Summary      Confirmar reserva
// @Description  Cambia el estado de la reserva a confirmed
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la reserva"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Reservation not found"
// @Router       /reservations/{id}/confirm [patch]
// @Security     Bearer
func ConfirmReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation models.Reservation

	if err := config.DB.First(&reservation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	now := time.Now()
	reservation.Status = "confirmed"
	reservation.ConfirmedAt = &now

	if err := config.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation confirmed successfully",
		"data":    reservation,
	})
}

// CancelReservation godoc
// @Summary      Cancelar reserva
// @Description  Cambia el estado de la reserva a cancelled
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la reserva"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Reservation not found"
// @Router       /reservations/{id}/cancel [patch]
// @Security     Bearer
func CancelReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation models.Reservation

	if err := config.DB.First(&reservation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	reservation.Status = "cancelled"

	if err := config.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation cancelled successfully",
		"data":    reservation,
	})
}

// SeatReservation godoc
// @Summary      Sentar reserva
// @Description  Cambia el estado de la reserva a seated y asigna mesa
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        id        path  int                      true  "ID de la reserva"
// @Param        request   body  map[string]interface{}   true  "table_id: ID de la mesa"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Reservation not found"
// @Router       /reservations/{id}/seat [patch]
// @Security     Bearer
func SeatReservation(c *gin.Context) {
	id := c.Param("id")
	var reservation models.Reservation

	if err := config.DB.First(&reservation, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	var request struct {
		TableID uint `json:"table_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()
	reservation.Status = "seated"
	reservation.TableID = &request.TableID
	reservation.SeatedAt = &now

	if err := config.DB.Save(&reservation).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to seat reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Reservation seated successfully",
		"data":    reservation,
	})
}

// GetReservationsByDate godoc
// @Summary      Obtener reservas por fecha
// @Description  Obtiene todas las reservas de una fecha específica
// @Tags         reservations
// @Accept       json
// @Produce      json
// @Param        date  path  string  true  "Fecha (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}  "data: array de reservations"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /reservations/date/{date} [get]
// @Security     Bearer
func GetReservationsByDate(c *gin.Context) {
	date := c.Param("date")
	var reservations []models.Reservation

	if err := config.DB.Where("reservation_date = ?", date).
		Preload("Company").Preload("Partner").Preload("Table").
		Find(&reservations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reservations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": reservations})
}
