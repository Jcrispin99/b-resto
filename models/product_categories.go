package models

import "gorm.io/gorm"

// ProductCategory - Categorización para menú/productos vendibles
type ProductCategory struct {
	gorm.Model
	ParentID *uint  `json:"parent_id" gorm:"default:null"`
	Type     string `json:"type" gorm:"size:50;default:'menu';not null"`
	Name     string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	FullName string `json:"full_name" gorm:"size:255"`
	IsActive bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Parent   *ProductCategory  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []ProductCategory `json:"children,omitempty" gorm:"foreignKey:ParentID"`
	Products []ProductTemplate `json:"products,omitempty" gorm:"foreignKey:CategoryID"`
}

func (ProductCategory) TableName() string {
	return "product_categories"
}
