package models

import "time"

// ProductAttributeValue - Valores específicos de un atributo (ej: "Pequeño", "Mediano", "Grande")
type ProductAttributeValue struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	AttributeID uint      `json:"attribute_id" gorm:"not null"`
	Value       string    `json:"value" gorm:"size:100;not null" binding:"required,min=1,max=100"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relaciones
	Attribute *ProductAttribute `json:"attribute,omitempty" gorm:"foreignKey:AttributeID"`
	Products  []ProductProduct  `json:"products,omitempty" gorm:"many2many:attribute_value_product;"`
}

func (ProductAttributeValue) TableName() string {
	return "product_attribute_values"
}
