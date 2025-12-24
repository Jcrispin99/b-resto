package models

import (
	"time"

	"gorm.io/gorm"
)

// Combo - Paquetes de productos con precio especial
type Combo struct {
	ID                 uint           `json:"id" gorm:"primaryKey"`
	CategoryID         uint           `json:"category_id" gorm:"not null"` // Para POS y menú
	Name               string         `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	Description        string         `json:"description" gorm:"type:text"`
	Price              float64        `json:"price" gorm:"type:decimal(10,2);not null"`
	RegularPrice       float64        `json:"regular_price" gorm:"type:decimal(10,2);not null"`
	DiscountPercentage float64        `json:"discount_percentage" gorm:"type:decimal(5,2);default:0;not null"`
	Image              string         `json:"image" gorm:"size:500"`
	StartDate          time.Time      `json:"start_date" gorm:"type:date;not null"`
	EndDate            *time.Time     `json:"end_date" gorm:"type:date"`
	IsActive           bool           `json:"is_active" gorm:"default:true;not null"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Relaciones
	Category *ProductCategory `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Items    []ComboItem      `json:"items,omitempty" gorm:"foreignKey:ComboID"`
}

func (Combo) TableName() string {
	return "combos"
}
