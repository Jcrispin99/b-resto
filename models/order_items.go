package models

import "gorm.io/gorm"

// OrderItem - Items de una orden (productos vendidos)
type OrderItem struct {
	gorm.Model
	OrderID       uint    `json:"order_id" gorm:"not null"`
	ProductID     uint    `json:"product_id" gorm:"not null"` // FK a product_product (variante)
	Quantity      float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`
	PriceUnit     float64 `json:"price_unit" gorm:"type:decimal(10,2);not null"` // Precio histórico
	PriceSubtotal float64 `json:"price_subtotal" gorm:"type:decimal(10,2);not null"`
	ProductNotes  string  `json:"product_notes" gorm:"type:text"` // "Sin cebolla", "Extra queso"

	// Relaciones
	Order   *Order          `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	Product *ProductProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"` // ✅ Directo a product_product
}

func (OrderItem) TableName() string {
	return "order_items"
}
