package models

import "gorm.io/gorm"

type InventoryCategory struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	ParentID uint   `json:"parent_id" gorm:"default:null"`
	FullName string `json:"full_name" gorm:"size:255;default:null"`
	IsActive bool   `json:"is_active" gorm:"default:true;not null"`
}

func (InventoryCategory) TableName() string {
	return "inventory_categories"
}
