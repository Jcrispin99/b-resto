package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetKitchenTickets godoc
// @Summary      Listar tickets de cocina
// @Description  Obtiene lista de tickets de cocina con filtros
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        status      query  string  false  "Filtrar por estado"  Enums(pending, preparing, ready, delivered)
// @Param        station_id  query  int     false  "Filtrar por estación"
// @Success      200  {object}  map[string]interface{}  "data: array de kitchen tickets"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /kitchen-tickets [get]
// @Security     Bearer
func GetKitchenTickets(c *gin.Context) {
	var tickets []models.KitchenTicket

	query := config.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if stationID := c.Query("station_id"); stationID != "" {
		query = query.Where("kitchen_station_id = ?", stationID)
	}

	if err := query.Preload("Order").Preload("KitchenStation").Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch kitchen tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}

// GetKitchenTicket godoc
// @Summary      Obtener ticket de cocina
// @Description  Obtiene un ticket de cocina por ID
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del ticket"
// @Success      200  {object}  map[string]interface{}  "data: kitchen ticket"
// @Failure      404  {object}  map[string]string       "error: Kitchen ticket not found"
// @Router       /kitchen-tickets/{id} [get]
// @Security     Bearer
func GetKitchenTicket(c *gin.Context) {
	id := c.Param("id")
	var ticket models.KitchenTicket

	if err := config.DB.Preload("Order").Preload("KitchenStation").First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen ticket not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": ticket})
}

// UpdateTicketStatus godoc
// @Summary      Actualizar estado de ticket
// @Description  Cambia el estado de un ticket de cocina
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        id      path  int                      true  "ID del ticket"
// @Param        status  body  map[string]interface{}   true  "status: nuevo estado"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Kitchen ticket not found"
// @Router       /kitchen-tickets/{id}/status [patch]
// @Security     Bearer
func UpdateTicketStatus(c *gin.Context) {
	id := c.Param("id")
	var ticket models.KitchenTicket

	if err := config.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen ticket not found"})
		return
	}

	var request struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ticket.State = request.Status

	if err := config.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket status updated successfully",
		"data":    ticket,
	})
}

// MarkTicketPreparing godoc
// @Summary      Marcar ticket como preparando
// @Description  Cambia el estado del ticket a preparing
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del ticket"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Kitchen ticket not found"
// @Router       /kitchen-tickets/{id}/preparing [patch]
// @Security     Bearer
func MarkTicketPreparing(c *gin.Context) {
	id := c.Param("id")
	var ticket models.KitchenTicket

	if err := config.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen ticket not found"})
		return
	}

	ticket.State = "preparing"

	if err := config.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket marked as preparing",
		"data":    ticket,
	})
}

// MarkTicketReady godoc
// @Summary      Marcar ticket como listo
// @Description  Cambia el estado del ticket a ready
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del ticket"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Kitchen ticket not found"
// @Router       /kitchen-tickets/{id}/ready [patch]
// @Security     Bearer
func MarkTicketReady(c *gin.Context) {
	id := c.Param("id")
	var ticket models.KitchenTicket

	if err := config.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen ticket not found"})
		return
	}

	ticket.State = "ready"

	if err := config.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket marked as ready",
		"data":    ticket,
	})
}

// MarkTicketDelivered godoc
// @Summary      Marcar ticket como entregado
// @Description  Cambia el estado del ticket a delivered
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del ticket"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Kitchen ticket not found"
// @Router       /kitchen-tickets/{id}/delivered [patch]
// @Security     Bearer
func MarkTicketDelivered(c *gin.Context) {
	id := c.Param("id")
	var ticket models.KitchenTicket

	if err := config.DB.First(&ticket, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Kitchen ticket not found"})
		return
	}

	ticket.State = "delivered"

	if err := config.DB.Save(&ticket).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update ticket"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket marked as delivered",
		"data":    ticket,
	})
}

// GetTicketsByStation godoc
// @Summary      Obtener tickets por estación
// @Description  Obtiene todos los tickets de una estación de cocina específica
// @Tags         kitchen-tickets
// @Accept       json
// @Produce      json
// @Param        station_id  path  int  true  "ID de la estación"
// @Success      200  {object}  map[string]interface{}  "data: array de kitchen tickets"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /kitchen-tickets/station/{station_id} [get]
// @Security     Bearer
func GetTicketsByStation(c *gin.Context) {
	stationID := c.Param("station_id")
	var tickets []models.KitchenTicket

	if err := config.DB.Where("kitchen_station_id = ?", stationID).
		Preload("Order").
		Find(&tickets).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch tickets"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": tickets})
}
