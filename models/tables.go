package models

import "gorm.io/gorm"

// Table - Mesas individuales del restaurante
type Table struct {
	gorm.Model
	CompanyID uint   `json:"company_id" gorm:"not null"`
	AreaID    uint   `json:"area_id" gorm:"not null"`
	Number    string `json:"number" gorm:"size:50;not null" binding:"required"`
	Capacity  int    `json:"capacity" gorm:"not null" binding:"required,min=1"`
	QRCode    string `json:"qr_code" gorm:"size:500"`
	IsActive  bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Company *Company   `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Area    *TableArea `json:"area,omitempty" gorm:"foreignKey:AreaID"`
	Orders  []Order    `json:"orders,omitempty" gorm:"foreignKey:TableID"`
}

func (Table) TableName() string {
	return "tables"
}
