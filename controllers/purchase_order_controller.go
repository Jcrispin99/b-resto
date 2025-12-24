package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"b-resto/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetPurchaseOrders godoc
// @Summary      Listar órdenes de compra
// @Description  Obtiene lista de todas las órdenes de compra
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        status      query  string  false  "Filtrar por estado"  Enums(draft, sent, received, cancelled)
// @Param        partner_id  query  int     false  "Filtrar por proveedor"
// @Success      200  {object}  map[string]interface{}  "data: array de purchase orders"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /purchase-orders [get]
// @Security     Bearer
func GetPurchaseOrders(c *gin.Context) {
	var orders []models.PurchaseOrder

	query := config.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if partnerID := c.Query("partner_id"); partnerID != "" {
		query = query.Where("partner_id = ?", partnerID)
	}

	if err := query.Preload("Partner").Preload("Warehouse").Preload("Items").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch purchase orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetPurchaseOrder godoc
// @Summary      Obtener orden de compra
// @Description  Obtiene una orden de compra por ID con todos sus detalles
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden de compra"
// @Success      200  {object}  map[string]interface{}  "data: purchase order completa"
// @Failure      404  {object}  map[string]string       "error: Purchase order not found"
// @Router       /purchase-orders/{id} [get]
// @Security     Bearer
func GetPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.PurchaseOrder

	if err := config.DB.
		Preload("Partner").
		Preload("Warehouse").
		Preload("Company").
		Preload("Items").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// CreatePurchaseOrder godoc
// @Summary      Crear orden de compra
// @Description  Crea una nueva orden de compra
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        order  body  models.PurchaseOrder  true  "Datos de la orden"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /purchase-orders [post]
// @Security     Bearer
func CreatePurchaseOrder(c *gin.Context) {
	var order models.PurchaseOrder

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Estado inicial
	if order.Status == "" {
		order.Status = "draft"
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Purchase order created successfully",
		"data":    order,
	})
}

// UpdatePurchaseOrder godoc
// @Summary      Actualizar orden de compra
// @Description  Actualiza los datos de una orden de compra existente
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        id     path  int                     true  "ID de la orden"
// @Param        order  body  models.PurchaseOrder    true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Purchase order not found"
// @Router       /purchase-orders/{id} [put]
// @Security     Bearer
func UpdatePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.PurchaseOrder

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	var updateData models.PurchaseOrder
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&order).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update purchase order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order updated successfully",
		"data":    order,
	})
}

// SendPurchaseOrder godoc
// @Summary      Enviar orden de compra
// @Description  Cambia el estado de la orden a sent
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Purchase order not found"
// @Router       /purchase-orders/{id}/send [patch]
// @Security     Bearer
func SendPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.PurchaseOrder

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	order.Status = "sent"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send purchase order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order sent successfully",
		"data":    order,
	})
}

// ReceivePurchaseOrder godoc
// @Summary      Recibir orden de compra
// @Description  Cambia el estado de la orden a received y registra entrada en inventario
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Purchase order not found"
// @Router       /purchase-orders/{id}/receive [patch]
// @Security     Bearer
func ReceivePurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.PurchaseOrder

	// Cargar orden de compra con items
	if err := config.DB.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	// Validar que no esté ya recibida
	if order.Status == "received" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Purchase order already received"})
		return
	}

	// Registrar entrada en inventario (Kardex)
	inventoryService := services.NewInventoryService()

	// Usar el warehouse_id de la orden de compra
	if err := inventoryService.RegisterPurchase(order.ID, order.Items, order.WarehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Inventory error: %s", err.Error())})
		return
	}

	// Actualizar estado de la orden
	order.Status = "received"
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive purchase order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order received successfully and inventory updated",
		"data":    order,
	})
}

// CancelPurchaseOrder godoc
// @Summary      Cancelar orden de compra
// @Description  Cambia el estado de la orden a cancelled
// @Tags         purchase-orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Purchase order not found"
// @Router       /purchase-orders/{id}/cancel [patch]
// @Security     Bearer
func CancelPurchaseOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.PurchaseOrder

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Purchase order not found"})
		return
	}

	order.Status = "cancelled"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel purchase order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Purchase order cancelled successfully",
		"data":    order,
	})
}
