package services

import (
	"b-resto/config"
	"b-resto/models"
	"errors"
	"fmt"
)

// InventoryService maneja la lógica de Kardex (movimientos de inventario)
type InventoryService struct{}

// NewInventoryService crea una nueva instancia del servicio
func NewInventoryService() *InventoryService {
	return &InventoryService{}
}

// RegisterSale registra salida de inventario por venta
func (s *InventoryService) RegisterSale(orderID uint, items []models.OrderItem, warehouseID uint) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range items {
		// Obtener último saldo del producto en este almacén
		var lastKardex models.Inventory
		result := tx.Where("product_id = ? AND warehouse_id = ?", item.ProductID, warehouseID).
			Order("id desc").
			First(&lastKardex)

		// Balance anterior (si no existe, es 0)
		previousBalance := float64(0)
		if result.Error == nil {
			previousBalance = lastKardex.QuantityBalance
		}

		// Validar stock suficiente
		if previousBalance < item.Quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient stock for product %d: available %.2f, required %.2f",
				item.ProductID, previousBalance, item.Quantity)
		}

		// Crear movimiento de SALIDA
		kardex := models.Inventory{
			ProductID:       item.ProductID,
			WarehouseID:     warehouseID,
			OrderID:         &orderID,
			Detail:          fmt.Sprintf("Venta - Order #%d", orderID),
			QuantityOut:     item.Quantity,
			QuantityBalance: previousBalance - item.Quantity,
		}

		if err := tx.Create(&kardex).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create kardex entry: %w", err)
		}
	}

	return tx.Commit().Error
}

// RegisterPurchase registra entrada de inventario por compra
func (s *InventoryService) RegisterPurchase(purchaseOrderID uint, items []models.PurchaseOrderItem, warehouseID uint) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range items {
		// Obtener último saldo
		var lastKardex models.Inventory
		result := tx.Where("product_id = ? AND warehouse_id = ?", item.ProductID, warehouseID).
			Order("id desc").
			First(&lastKardex)

		previousBalance := float64(0)
		if result.Error == nil {
			previousBalance = lastKardex.QuantityBalance
		}

		// Crear movimiento de ENTRADA
		kardex := models.Inventory{
			ProductID:       item.ProductID,
			WarehouseID:     warehouseID,
			PurchaseOrderID: &purchaseOrderID,
			Detail:          fmt.Sprintf("Compra - Purchase Order #%d", purchaseOrderID),
			QuantityIn:      item.Quantity,
			CostIn:          item.UnitPrice,
			TotalIn:         item.Quantity * item.UnitPrice,
			QuantityBalance: previousBalance + item.Quantity,
		}

		if err := tx.Create(&kardex).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create kardex entry: %w", err)
		}
	}

	return tx.Commit().Error
}

// RegisterTransfer registra transferencia entre almacenes (salida + entrada)
func (s *InventoryService) RegisterTransfer(transferID uint, fromWarehouseID, toWarehouseID uint, items []models.StockTransferItem) error {
	tx := config.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for _, item := range items {
		// 1. SALIDA del almacén origen
		var lastKardexFrom models.Inventory
		result := tx.Where("product_id = ? AND warehouse_id = ?", item.ProductID, fromWarehouseID).
			Order("id desc").
			First(&lastKardexFrom)

		previousBalanceFrom := float64(0)
		if result.Error == nil {
			previousBalanceFrom = lastKardexFrom.QuantityBalance
		}

		// Validar stock en origen
		if previousBalanceFrom < item.Quantity {
			tx.Rollback()
			return fmt.Errorf("insufficient stock in source warehouse for product %d", item.ProductID)
		}

		kardexOut := models.Inventory{
			ProductID:       item.ProductID,
			WarehouseID:     fromWarehouseID,
			StockTransferID: &transferID,
			Detail:          fmt.Sprintf("Transferencia salida - Transfer #%d", transferID),
			QuantityOut:     item.Quantity,
			QuantityBalance: previousBalanceFrom - item.Quantity,
		}

		if err := tx.Create(&kardexOut).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create kardex out entry: %w", err)
		}

		// 2. ENTRADA al almacén destino
		var lastKardexTo models.Inventory
		result = tx.Where("product_id = ? AND warehouse_id = ?", item.ProductID, toWarehouseID).
			Order("id desc").
			First(&lastKardexTo)

		previousBalanceTo := float64(0)
		if result.Error == nil {
			previousBalanceTo = lastKardexTo.QuantityBalance
		}

		kardexIn := models.Inventory{
			ProductID:       item.ProductID,
			WarehouseID:     toWarehouseID,
			StockTransferID: &transferID,
			Detail:          fmt.Sprintf("Transferencia entrada - Transfer #%d", transferID),
			QuantityIn:      item.Quantity,
			QuantityBalance: previousBalanceTo + item.Quantity,
		}

		if err := tx.Create(&kardexIn).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to create kardex in entry: %w", err)
		}
	}

	return tx.Commit().Error
}

// ValidateStock verifica si hay stock suficiente antes de una venta
func (s *InventoryService) ValidateStock(productID, warehouseID uint, requiredQty float64) (bool, float64, error) {
	var lastKardex models.Inventory
	result := config.DB.Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		Order("id desc").
		First(&lastKardex)

	if result.Error != nil {
		// No hay movimientos = stock 0
		return false, 0, nil
	}

	available := lastKardex.QuantityBalance
	hasStock := available >= requiredQty

	return hasStock, available, nil
}

// GetCurrentStock obtiene el stock actual de un producto en un almacén
func (s *InventoryService) GetCurrentStock(productID, warehouseID uint) (float64, error) {
	var lastKardex models.Inventory
	result := config.DB.Where("product_id = ? AND warehouse_id = ?", productID, warehouseID).
		Order("id desc").
		First(&lastKardex)

	if result.Error != nil {
		return 0, errors.New("no inventory records found")
	}

	return lastKardex.QuantityBalance, nil
}
