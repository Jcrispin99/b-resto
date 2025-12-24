package models

import "gorm.io/gorm"

// TableArea - Áreas de mesas (Salón, Terraza, VIP, etc.)
type TableArea struct {
	gorm.Model
	CompanyID   uint   `json:"company_id" gorm:"not null"`
	Name        string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	Description string `json:"description" gorm:"type:text"`
	Order       int    `json:"order" gorm:"default:0;not null"`
	IsActive    bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Tables  []Table  `json:"tables,omitempty" gorm:"foreignKey:AreaID"`
}

func (TableArea) TableName() string {
	return "table_areas"
}
