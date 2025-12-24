package models

import "gorm.io/gorm"

// PurchaseOrderItem - Items de una orden de compra
type PurchaseOrderItem struct {
	gorm.Model
	PurchaseOrderID uint    `json:"purchase_order_id" gorm:"not null"`
	ProductID       uint    `json:"product_id" gorm:"not null"` // FK a product_product
	Quantity        float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`
	UnitPrice       float64 `json:"unit_price" gorm:"type:decimal(10,2);not null"`
	Subtotal        float64 `json:"subtotal" gorm:"type:decimal(10,2);not null"`

	// Relaciones
	PurchaseOrder *PurchaseOrder  `json:"purchase_order,omitempty" gorm:"foreignKey:PurchaseOrderID"`
	Product       *ProductProduct `json:"product,omitempty" gorm:"foreignKey:ProductID"`
}

func (PurchaseOrderItem) TableName() string {
	return "purchase_order_items"
}
