package models

import "gorm.io/gorm"

// KitchenStation - Estaciones de cocina donde se preparan los productos
type KitchenStation struct {
	gorm.Model
	CompanyID   uint   `json:"company_id" gorm:"not null"`
	Name        string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	Description string `json:"description" gorm:"size:500"`
	PrinterIP   string `json:"printer_ip" gorm:"size:50;column:printer_ip"`
	Order       int    `json:"order" gorm:"default:0;not null"`
	IsActive    bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Company  *Company          `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Products []ProductTemplate `json:"products,omitempty" gorm:"foreignKey:KitchenStationID"`
}

func (KitchenStation) TableName() string {
	return "kitchen_stations"
}
