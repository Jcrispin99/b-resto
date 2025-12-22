package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPaymentMethods obtiene todos los métodos de pago
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

// GetPaymentMethod obtiene un método de pago por ID
func GetPaymentMethod(c *gin.Context) {
	id := c.Param("id")
	var paymentMethod models.PaymentMethod

	if err := config.DB.First(&paymentMethod, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment method not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": paymentMethod})
}

// CreatePaymentMethod crea un nuevo método de pago
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

// UpdatePaymentMethod actualiza un método de pago existente
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

// DeletePaymentMethod elimina un método de pago (soft delete)
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

// TogglePaymentMethodStatus activa/desactiva un método de pago
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
