package models

import "time"

// ComboItem - Productos que componen un combo
type ComboItem struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	ComboID           uint      `json:"combo_id" gorm:"not null"`
	ProductTemplateID uint      `json:"product_template_id" gorm:"not null"`
	Quantity          int       `json:"quantity" gorm:"default:1;not null"`
	AllowSubstitution bool      `json:"allow_substitution" gorm:"default:false;not null"`
	CreatedAt         time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relaciones
	Combo           *Combo           `json:"combo,omitempty" gorm:"foreignKey:ComboID"`
	ProductTemplate *ProductTemplate `json:"product_template,omitempty" gorm:"foreignKey:ProductTemplateID"`
}

func (ComboItem) TableName() string {
	return "combo_items"
}
