package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOrderPayments godoc
// @Summary      Listar pagos de orden
// @Description  Obtiene lista de pagos de una orden específica
// @Tags         order-payments
// @Accept       json
// @Produce      json
// @Param        order_id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "data: array de payments"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /orders/{order_id}/payments [get]
// @Security     Bearer
func GetOrderPayments(c *gin.Context) {
	orderID := c.Param("id")
	var payments []models.OrderPayment

	if err := config.DB.Where("order_id = ?", orderID).
		Preload("PaymentMethod").
		Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payments"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": payments})
}

// CreateOrderPayment godoc
// @Summary      Crear pago de orden
// @Description  Registra un nuevo pago para una orden
// @Tags         order-payments
// @Accept       json
// @Produce      json
// @Param        order_id  path  int                  true  "ID de la orden"
// @Param        payment   body  models.OrderPayment  true  "Datos del pago"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /orders/{order_id}/payments [post]
// @Security     Bearer
func CreateOrderPayment(c *gin.Context) {
	orderID := c.Param("id")
	var payment models.OrderPayment

	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Asignar el OrderID del path
	var orderIDUint uint
	if _, err := fmt.Sscanf(orderID, "%d", &orderIDUint); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid order ID"})
		return
	}
	payment.OrderID = orderIDUint

	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment created successfully",
		"data":    payment,
	})
}

// DeleteOrderPayment godoc
// @Summary      Eliminar pago de orden
// @Description  Elimina un pago de una orden
// @Tags         order-payments
// @Accept       json
// @Produce      json
// @Param        order_id    path  int  true  "ID de la orden"
// @Param        payment_id  path  int  true  "ID del pago"
// @Success      200  {object}  map[string]string  "message: Payment deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Payment not found"
// @Router       /orders/{order_id}/payments/{payment_id} [delete]
// @Security     Bearer
func DeleteOrderPayment(c *gin.Context) {
	paymentID := c.Param("payment_id")
	var payment models.OrderPayment

	if err := config.DB.First(&payment, paymentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment not found"})
		return
	}

	if err := config.DB.Delete(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment deleted successfully"})
}
