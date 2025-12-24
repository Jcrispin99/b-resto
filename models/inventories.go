package models

import (
	"time"

	"gorm.io/gorm"
)

// Inventory - Movimientos de inventario (Kardex)
// ✅ SIN polimorfismo - usando campos específicos por tipo de origen
type Inventory struct {
	ID          uint `json:"id" gorm:"primaryKey"`
	ProductID   uint `json:"product_id" gorm:"not null"`
	WarehouseID uint `json:"warehouse_id" gorm:"not null"`

	// ✅ Origen del movimiento (solo uno debe estar set)
	OrderID         *uint `json:"order_id"`          // Si es por venta
	PurchaseOrderID *uint `json:"purchase_order_id"` // Si es por compra
	StockTransferID *uint `json:"stock_transfer_id"` // Si es por transferencia

	Detail string `json:"detail" gorm:"size:500"` // Descripción (ajustes manuales)

	// Entradas
	QuantityIn float64 `json:"quantity_in" gorm:"type:decimal(10,4);default:0;not null"`
	CostIn     float64 `json:"cost_in" gorm:"type:decimal(10,2);default:0;not null"`
	TotalIn    float64 `json:"total_in" gorm:"type:decimal(10,2);default:0;not null"`

	// Salidas
	QuantityOut float64 `json:"quantity_out" gorm:"type:decimal(10,4);default:0;not null"`
	CostOut     float64 `json:"cost_out" gorm:"type:decimal(10,2);default:0;not null"`
	TotalOut    float64 `json:"total_out" gorm:"type:decimal(10,2);default:0;not null"`

	// Balances (acumulados)
	QuantityBalance float64 `json:"quantity_balance" gorm:"type:decimal(10,4);default:0;not null"`
	CostBalance     float64 `json:"cost_balance" gorm:"type:decimal(10,2);default:0;not null"`
	TotalBalance    float64 `json:"total_balance" gorm:"type:decimal(10,2);default:0;not null"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relaciones
	Product       *ProductProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
	Warehouse     *Warehouse      `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Order         *Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	PurchaseOrder *PurchaseOrder  `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	StockTransfer *StockTransfer  `json:"stock_transfer,omitempty" gorm:"foreignKey:StockTransferID"`
}

func (Inventory) TableName() string {
	return "inventories"
}
