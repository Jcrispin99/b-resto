package models

import (
	"gorm.io/gorm"
)

// ProductTemplate - Maestro de productos, contiene la información común
type ProductTemplate struct {
	gorm.Model
	InventoryCategoryID *uint   `json:"inventory_category_id" gorm:"default:null"`
	CategoryID          uint    `json:"category_id" gorm:"not null"`
	UnitID              uint    `json:"unit_id" gorm:"not null"`
	KitchenStationID    *uint   `json:"kitchen_station_id" gorm:"default:null"`
	Name                string  `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	Description         string  `json:"description" gorm:"type:text"`
	InternalReference   string  `json:"internal_reference" gorm:"size:100"`
	Barcode             string  `json:"barcode" gorm:"size:100"`
	ProductType         string  `json:"product_type" gorm:"size:50;default:'storable';not null"` // storable, service, consumable
	CanBeSold           bool    `json:"can_be_sold" gorm:"default:false;not null"`
	CanBePurchased      bool    `json:"can_be_purchased" gorm:"default:true;not null"`
	CanBeStocked        bool    `json:"can_be_stocked" gorm:"default:true;not null"`
	SalePrice           float64 `json:"sale_price" gorm:"type:decimal(10,2);default:0;not null"`
	IsActive            bool    `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	InventoryCategory *InventoryCategory `json:"inventory_category,omitempty" gorm:"foreignKey:InventoryCategoryID"`
	Category          *ProductCategory   `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Unit              *Unit              `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
	KitchenStation    *KitchenStation    `json:"kitchen_station,omitempty" gorm:"foreignKey:KitchenStationID"`
	Variants          []ProductProduct   `json:"variants,omitempty" gorm:"foreignKey:TemplateID"`
}

func (ProductTemplate) TableName() string {
	return "product_template"
}
