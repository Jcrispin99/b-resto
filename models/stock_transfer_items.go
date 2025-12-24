package models

import "gorm.io/gorm"

// StockTransferItem - Items de una transferencia de stock
type StockTransferItem struct {
	gorm.Model
	StockTransferID uint    `json:"stock_transfer_id" gorm:"not null"`
	ProductID       uint    `json:"product_id" gorm:"not null"` // FK a product_product
	Quantity        float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`
	Cost            float64 `json:"cost" gorm:"type:decimal(10,2);not null"`

	// Relaciones
	StockTransfer *StockTransfer  `json:"stock_transfer,omitempty" gorm:"foreignKey:StockTransferID"`
	Product       *ProductProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (StockTransferItem) TableName() string {
	return "stock_transfer_items"
}
