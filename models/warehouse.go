package models

import "gorm.io/gorm"

// Warehouse - Almacenes/bodegas donde se guarda inventario
type Warehouse struct {
	gorm.Model
	CompanyID uint   `json:"company_id" gorm:"not null"`
	Code      string `json:"code" gorm:"size:50;not null;uniqueIndex:idx_warehouse_code"`
	Name      string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	IsActive  bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
}

func (Warehouse) TableName() string {
	return "warehouses"
}
