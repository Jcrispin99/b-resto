package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetInventories godoc
// @Summary      Listar inventarios (kardex)
// @Description  Obtiene el inventario/stock actual de productos por almacén
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        warehouse_id  query  int     false  "Filtrar por almacén"
// @Param        product_id    query  int     false  "Filtrar por producto"
// @Success      200  {object}  map[string]interface{}  "data: array de inventories"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /inventories [get]
// @Security     Bearer
func GetInventories(c *gin.Context) {
	var inventories []models.Inventory

	query := config.DB
	if warehouseID := c.Query("warehouse_id"); warehouseID != "" {
		query = query.Where("warehouse_id = ?", warehouseID)
	}
	if productID := c.Query("product_id"); productID != "" {
		query = query.Where("product_id = ?", productID)
	}

	if err := query.Preload("Warehouse").Preload("Product").Find(&inventories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inventories"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventories})
}

// GetInventory godoc
// @Summary      Obtener inventario
// @Description  Obtiene el inventario de un producto en un almacén específico
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID del inventario"
// @Success      200  {object}  map[string]interface{}  "data: inventory"
// @Failure      404  {object}  map[string]string       "error: Inventory not found"
// @Router       /inventories/{id} [get]
// @Security     Bearer
func GetInventory(c *gin.Context) {
	id := c.Param("id")
	var inventory models.Inventory

	if err := config.DB.Preload("Warehouse").Preload("Product").First(&inventory, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventory})
}

// GetInventoryByWarehouseAndProduct godoc
// @Summary      Obtener stock específico
// @Description  Obtiene el stock de un producto en un almacén
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        warehouse_id  path  int  true  "ID del almacén"
// @Param        product_id    path  int  true  "ID del producto"
// @Success      200  {object}  map[string]interface{}  "data: inventory"
// @Failure      404  {object}  map[string]string       "error: Inventory not found"
// @Router       /inventories/warehouse/{warehouse_id}/product/{product_id} [get]
// @Security     Bearer
func GetInventoryByWarehouseAndProduct(c *gin.Context) {
	warehouseID := c.Param("warehouse_id")
	productID := c.Param("product_id")
	var inventory models.Inventory

	if err := config.DB.
		Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).
		Preload("Warehouse").
		Preload("Product").
		First(&inventory).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Inventory not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventory})
}

// AdjustInventory godoc
// @Summary      Ajustar inventario
// @Description  Realiza un ajuste manual de inventario (entrada/salida)
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Param        adjustment  body  map[string]interface{}  true  "warehouse_id, product_id, quantity, type (in/out), reason"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /inventories/adjust [post]
// @Security     Bearer
func AdjustInventory(c *gin.Context) {
	var request struct {
		WarehouseID uint    `json:"warehouse_id" binding:"required"`
		ProductID   uint    `json:"product_id" binding:"required"`
		Quantity    float64 `json:"quantity" binding:"required"`
		Type        string  `json:"type" binding:"required,oneof=in out"` // in/out
		Reason      string  `json:"reason"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Buscar o crear inventario
	var inventory models.Inventory
	result := config.DB.Where("warehouse_id = ? AND product_id = ?", request.WarehouseID, request.ProductID).First(&inventory)

	if result.Error != nil {
		// Crear nuevo registro de inventario
		inventory = models.Inventory{
			WarehouseID:     request.WarehouseID,
			ProductID:       request.ProductID,
			QuantityBalance: 0,
		}
	}

	// Ajustar cantidad
	if request.Type == "in" {
		inventory.QuantityBalance += request.Quantity
	} else {
		inventory.QuantityBalance -= request.Quantity
		if inventory.QuantityBalance < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient stock"})
			return
		}
	}

	// Guardar o actualizar
	if result.Error != nil {
		if err := config.DB.Create(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create inventory"})
			return
		}
	} else {
		if err := config.DB.Save(&inventory).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update inventory"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Inventory adjusted successfully",
		"data":    inventory,
	})
}

// GetLowStockProducts godoc
// @Summary      Productos con stock bajo
// @Description  Obtiene productos que están por debajo del stock mínimo
// @Tags         inventories
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "data: array de inventories con stock bajo"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /inventories/low-stock [get]
// @Security     Bearer
func GetLowStockProducts(c *gin.Context) {
	var inventories []models.Inventory

	// Productos donde quantity_on_hand <= min_quantity
	// Nota: El modelo Inventory no tiene min_quantity, ajustar según necesidad
	if err := config.DB.
		Where("quantity_balance < ?", 10). // Stock bajo si < 10
		Preload("Warehouse").
		Preload("Product").
		Find(&inventories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch low stock products"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": inventories})
}
