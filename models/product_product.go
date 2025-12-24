package models

import (
	"time"

	"gorm.io/gorm"
)

// ProductProduct - Variantes de productos (ej: Pizza Grande, Pizza Mediana)
type ProductProduct struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	TemplateID uint           `json:"template_id" gorm:"not null"`
	SKU        string         `json:"sku" gorm:"size:100;not null;uniqueIndex"`
	Barcode    string         `json:"barcode" gorm:"size:100"`
	SalePrice  *float64       `json:"sale_price" gorm:"type:decimal(10,2)"` // Sobrescribe el precio del template
	IsActive   bool           `json:"is_active" gorm:"default:true;not null"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relaciones
	Template        *ProductTemplate        `json:"template,omitempty" gorm:"foreignKey:TemplateID"`
	AttributeValues []ProductAttributeValue `json:"attribute_values,omitempty" gorm:"many2many:attribute_value_product;"`
}

func (ProductProduct) TableName() string {
	return "product_product"
}
