package models

import "time"

// ProductAttribute - Define tipos de atributos (ej: "Tamaño", "Temperatura")
type ProductAttribute struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name" gorm:"size:100;not null" binding:"required,min=2,max=100"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relaciones
	Values []ProductAttributeValue `json:"values,omitempty" gorm:"foreignKey:AttributeID"`
}

func (ProductAttribute) TableName() string {
	return "product_attributes"
}
