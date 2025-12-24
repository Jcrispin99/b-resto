package controllers

import (
	"b-resto/config"
	"b-resto/models"
	"b-resto/services"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetStockTransfers godoc
// @Summary      Listar transferencias de stock
// @Description  Obtiene lista de todas las transferencias entre almacenes
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        status  query  string  false  "Filtrar por estado"  Enums(draft, in_transit, received, cancelled)
// @Success      200  {object}  map[string]interface{}  "data: array de stock transfers"
// @Failure      500  {object}  map[string]string       "error: mensaje"
// @Router       /stock-transfers [get]
// @Security     Bearer
func GetStockTransfers(c *gin.Context) {
	var transfers []models.StockTransfer

	query := config.DB
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("Items").
		Find(&transfers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch stock transfers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transfers})
}

// GetStockTransfer godoc
// @Summary      Obtener transferencia de stock
// @Description  Obtiene una transferencia por ID con todos sus detalles
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la transferencia"
// @Success      200  {object}  map[string]interface{}  "data: stock transfer completa"
// @Failure      404  {object}  map[string]string       "error: Stock transfer not found"
// @Router       /stock-transfers/{id} [get]
// @Security     Bearer
func GetStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var transfer models.StockTransfer

	if err := config.DB.
		Preload("FromWarehouse").
		Preload("ToWarehouse").
		Preload("Items").
		First(&transfer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock transfer not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transfer})
}

// CreateStockTransfer godoc
// @Summary      Crear transferencia de stock
// @Description  Crea una nueva transferencia entre almacenes
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        transfer  body  models.StockTransfer  true  "Datos de la transferencia"
// @Success      201  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Router       /stock-transfers [post]
// @Security     Bearer
func CreateStockTransfer(c *gin.Context) {
	var transfer models.StockTransfer

	if err := c.ShouldBindJSON(&transfer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Estado inicial
	if transfer.Status == "" {
		transfer.Status = "draft"
	}

	// Validar que origen y destino sean diferentes
	if transfer.FromWarehouseID == transfer.ToWarehouseID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Source and destination warehouses must be different"})
		return
	}

	if err := config.DB.Create(&transfer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create stock transfer"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Stock transfer created successfully",
		"data":    transfer,
	})
}

// UpdateStockTransfer godoc
// @Summary      Actualizar transferencia de stock
// @Description  Actualiza los datos de una transferencia existente
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        id        path  int                    true  "ID de la transferencia"
// @Param        transfer  body  models.StockTransfer   true  "Datos actualizados"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      400  {object}  map[string]string       "error: validación"
// @Failure      404  {object}  map[string]string       "error: Stock transfer not found"
// @Router       /stock-transfers/{id} [put]
// @Security     Bearer
func UpdateStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var transfer models.StockTransfer

	if err := config.DB.First(&transfer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock transfer not found"})
		return
	}

	var updateData models.StockTransfer
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Model(&transfer).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update stock transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock transfer updated successfully",
		"data":    transfer,
	})
}

// SendStockTransfer godoc
// @Summary      Enviar transferencia
// @Description  Cambia el estado a in_transit (descuenta del origen)
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la transferencia"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Stock transfer not found"
// @Router       /stock-transfers/{id}/send [patch]
// @Security     Bearer
func SendStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var transfer models.StockTransfer

	if err := config.DB.First(&transfer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock transfer not found"})
		return
	}

	transfer.Status = "in_transit"
	// TODO: Descontar stock del almacén origen

	if err := config.DB.Save(&transfer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send stock transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock transfer sent successfully",
		"data":    transfer,
	})
}

// ReceiveStockTransfer godoc
// @Summary      Recibir transferencia
// @Description  Cambia el estado a received y registra movimientos en inventario (salida origen + entrada destino)
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la transferencia"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Stock transfer not found"
// @Router       /stock-transfers/{id}/receive [patch]
// @Security     Bearer
func ReceiveStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var transfer models.StockTransfer

	// Cargar transferencia con items
	if err := config.DB.Preload("Items").First(&transfer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock transfer not found"})
		return
	}

	// Validar que no esté ya recibida
	if transfer.Status == "received" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stock transfer already received"})
		return
	}

	// Registrar movimientos en inventario (Kardex): salida del origen + entrada al destino
	inventoryService := services.NewInventoryService()

	if err := inventoryService.RegisterTransfer(transfer.ID, transfer.FromWarehouseID, transfer.ToWarehouseID, transfer.Items); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Inventory error: %s", err.Error())})
		return
	}

	// Actualizar estado de la transferencia
	transfer.Status = "received"
	if err := config.DB.Save(&transfer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to receive stock transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock transfer received successfully and inventory updated",
		"data":    transfer,
	})
}

// CancelStockTransfer godoc
// @Summary      Cancelar transferencia
// @Description  Cambia el estado a cancelled
// @Tags         stock-transfers
// @Accept       json
// @Produce      json
// @Param        id  path  int  true  "ID de la transferencia"
// @Success      200  {object}  map[string]interface{}  "message y data"
// @Failure      404  {object}  map[string]string       "error: Stock transfer not found"
// @Router       /stock-transfers/{id}/cancel [patch]
// @Security     Bearer
func CancelStockTransfer(c *gin.Context) {
	id := c.Param("id")
	var transfer models.StockTransfer

	if err := config.DB.First(&transfer, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Stock transfer not found"})
		return
	}

	transfer.Status = "cancelled"

	if err := config.DB.Save(&transfer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel stock transfer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Stock transfer cancelled successfully",
		"data":    transfer,
	})
}
