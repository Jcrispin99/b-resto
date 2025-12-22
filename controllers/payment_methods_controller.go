package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaymentMethods godoc
// @Summary      Listar métodos de pago
// @Description  Obtiene lista de todos los métodos de pago disponibles
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        is_active  query  string  false  "Filtrar por estado activo"  Enums(true, false)
// @Success      200  {object}  map[string]interface{}  "data: array de payment methods"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /payment-methods [get]
// @Security     Bearer
func GetPaymentMethods(c *gin.Context) {
	var paymentMethods []models.PaymentMethod

	query := config.DB
	if isActive := c.Query("is_active"); isActive != "" {
		query = query.Where("is_active = ?", isActive == "true")
	}

	if err := query.Find(&paymentMethods).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment methods"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": paymentMethods})
}

// GetPaymentMethod godoc
// @Summary      Obtener método de pago
// @Description  Obtiene un método de pago por ID
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del método de pago"
// @Success      200  {object}  map[string]interface{}  "data: payment method"
// @Failure      404  {object}  map[string]string       "error: Payment method not found"
// @Router       /payment-methods/{id} [get]
// @Security     Bearer
func GetPaymentMethod(c *gin.Context) {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod

	if err := config.DB.First(&paymentMethod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": paymentMethod})
}

// CreatePaymentMethod godoc
// @Summary      Crear método de pago
// @Description  Crea un nuevo método de pago
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        paymentMethod  body  models.PaymentMethod  true  "Datos del método de pago"
// @Success      201  {object}  map[string]interface{}  "message y data creada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      409  {object}  map[string]string       "error: código duplicado"
// @Router       /payment-methods [post]
// @Security     Bearer
func CreatePaymentMethod(c *gin.Context) {
	var paymentMethod models.PaymentMethod

	// Gin valida automáticamente según los binding tags del modelo
	if err := c.ShouldBindJSON(&paymentMethod); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Solo validaciones de negocio (unicidad de código)
	var existing models.PaymentMethod
	if err := config.DB.Where("code = ?", paymentMethod.Code).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Payment method with this code already exists"})
		return
	}

	if err := config.DB.Create(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment method"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Payment method created successfully",
		"data":    paymentMethod,
	})
}

// UpdatePaymentMethod godoc
// @Summary      Actualizar método de pago
// @Description  Actualiza los datos de un método de pago existente
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        id             path  int                   true  "ID del método de pago"
// @Param        paymentMethod  body  models.PaymentMethod  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data actualizada"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Payment method not found"
// @Failure      409  {object}  map[string]string       "error: código duplicado"
// @Router       /payment-methods/{id} [put]
// @Security     Bearer
func UpdatePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod

	if err := config.DB.First(&paymentMethod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	var updateData models.PaymentMethod
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validar código único si se está cambiando
	if updateData.Code != "" && updateData.Code != paymentMethod.Code {
		var existing models.PaymentMethod
		if err := config.DB.Where("code = ? AND id != ?", updateData.Code, id).First(&existing).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Code already exists"})
			return
		}
	}

	if err := config.DB.Model(&paymentMethod).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Payment method updated successfully",
		"data":    paymentMethod,
	})
}

// DeletePaymentMethod godoc
// @Summary      Eliminar método de pago
// @Description  Elimina un método de pago (soft delete)
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del método de pago"
// @Success      200  {object}  map[string]string  "message: Payment method deleted successfully"
// @Failure      404  {object}  map[string]string  "error: Payment method not found"
// @Router       /payment-methods/{id} [delete]
// @Security     Bearer
func DeletePaymentMethod(c *gin.Context) {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod

	if err := config.DB.First(&paymentMethod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	if err := config.DB.Delete(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment method"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Payment method deleted successfully"})
}

// TogglePaymentMethodStatus godoc
// @Summary      Activar/Desactivar método de pago
// @Description  Cambia el estado is_active de un método de pago
// @Tags         payment-methods
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del método de pago"
// @Success      200  {object}  map[string]interface{}  "message y data con nuevo estado"
// @Failure      404  {object}  map[string]string       "error: Payment method not found"
// @Router       /payment-methods/{id}/toggle [patch]
// @Security     Bearer
func TogglePaymentMethodStatus(c *gin.Context) {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod

	if err := config.DB.First(&paymentMethod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	paymentMethod.IsActive = !paymentMethod.IsActive

	if err := config.DB.Save(&paymentMethod).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to toggle status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status toggled successfully",
		"data":    paymentMethod,
	})
}
