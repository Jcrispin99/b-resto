package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"b-resto/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetOrders godoc
// @Summary      Listar órdenes
// @Description  Obtiene lista de todas las órdenes con filtros
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        state     query  string  false  "Filtrar por estado"  Enums(draft, confirmed, done, cancelled)
// @Param        table_id  query  int     false  "Filtrar por mesa"
// @Param        date      query  string  false  "Filtrar por fecha (YYYY-MM-DD)"
// @Success      200  {object}  map[string]interface{}  "data: array de orders"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /orders [get]
// @Security     Bearer
func GetOrders(c *gin.Context) {
	var orders []models.Order

	query := config.DB
	if state := c.Query("state"); state != "" {
		query = query.Where("state = ?", state)
	}
	if tableID := c.Query("table_id"); tableID != "" {
		query = query.Where("table_id = ?", tableID)
	}
	if date := c.Query("date"); date != "" {
		query = query.Where("order_date = ?", date)
	}

	if err := query.Preload("Journal").Preload("User").Preload("Items").Preload("Payments").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}

// GetOrder godoc
// @Summary      Obtener orden
// @Description  Obtiene una orden por ID con todos sus detalles
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "data: order completa"
// @Failure      404  {object}  map[string]string       "error: Order not found"
// @Router       /orders/{id} [get]
// @Security     Bearer
func GetOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.
		Preload("Journal").
		Preload("User").
		Preload("Items").
		Preload("Payments").
		Preload("Tickets").
		First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": order})
}

// CreateOrder godoc
// @Summary      Crear orden
// @Description  Crea una nueva orden de venta
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        order  body  models.Order  true  "Datos de la orden"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /orders [post]
// @Security     Bearer
func CreateOrder(c *gin.Context) {
	var order models.Order

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Estado inicial
	if order.State == "" {
		order.State = "draft"
	}

	if err := config.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"data":    order,
	})
}

// UpdateOrder godoc
// @Summary      Actualizar orden
// @Description  Actualiza los datos de una orden existente
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id     path  int           true  "ID de la orden"
// @Param        order  body  models.Order  true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Order not found"
// @Router       /orders/{id} [put]
// @Security     Bearer
func UpdateOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	var updateData models.Order
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&order).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order updated successfully",
		"data":    order,
	})
}

// ConfirmOrder godoc
// @Summary      Confirmar orden
// @Description  Cambia el estado de la orden a confirmed
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Order not found"
// @Router       /orders/{id}/confirm [patch]
// @Security     Bearer
func ConfirmOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.State = "confirmed"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to confirm order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order confirmed successfully",
		"data":    order,
	})
}

// CancelOrder godoc
// @Summary      Cancelar orden
// @Description  Cambia el estado de la orden a cancelled
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Order not found"
// @Router       /orders/{id}/cancel [patch]
// @Security     Bearer
func CancelOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	if err := config.DB.First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	order.State = "cancelled"

	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
		"data":    order,
	})
}

// CompleteOrder godoc
// @Summary      Completar orden
// @Description  Cambia el estado de la orden a done y registra salida en inventario
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la orden"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Order not found"
// @Router       /orders/{id}/complete [patch]
// @Security     Bearer
func CompleteOrder(c *gin.Context) {
	id := c.Param("id")
	var order models.Order

	// Cargar orden con items
	if err := config.DB.Preload("Items").First(&order, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}

	// Validar que no esté ya completada
	if order.State == "done" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Order already completed"})
		return
	}

	// Registrar salida en inventario (Kardex)
	inventoryService := services.NewInventoryService()

	// TODO: Obtener warehouse_id del contexto o configuración
	// Por ahora usamos warehouse 4 (el que creamos en pruebas)
	warehouseID := uint(4)

	if err := inventoryService.RegisterSale(order.ID, order.Items, warehouseID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Inventory error: %s", err.Error())})
		return
	}

	// Actualizar estado de la orden
	order.State = "done"
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete order"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order completed successfully and inventory updated",
		"data":    order,
	})
}

// GetOrdersByTable godoc
// @Summary      Obtener órdenes por mesa
// @Description  Obtiene todas las órdenes de una mesa específica
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        table_id  path  int  true  "ID de la mesa"
// @Success      200  {object}  map[string]interface{}  "data: array de orders"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /orders/table/{table_id} [get]
// @Security     Bearer
func GetOrdersByTable(c *gin.Context) {
	tableID := c.Param("table_id")
	var orders []models.Order

	if err := config.DB.Where("table_id = ?", tableID).
		Preload("Items").Preload("Payments").
		Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch orders"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": orders})
}
