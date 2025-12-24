package models

import "gorm.io/gorm"

// InventoryCategory - Categorización para inventario (materias primas, insumos)
type InventoryCategory struct {
	gorm.Model
	ParentID *uint  `json:"parent_id" gorm:"default:null"`
	Name     string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	FullName string `json:"full_name" gorm:"size:255"`
	IsActive bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Parent   *InventoryCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []InventoryCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []ProductTemplate   `json:"products,omitempty" gorm:"foreignKey:InventoryCategoryID"`
}

func (InventoryCategory) TableName() string {
	return "inventory_categories"
}
