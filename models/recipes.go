package models

import "gorm.io/gorm"

// Recipe - Recetas/BOM (Bill of Materials) - Ingredientes por producto
type Recipe struct {
	gorm.Model
	ProductTemplateID uint    `json:"product_template_id" gorm:"not null"` // Producto final
	IngredientID      uint    `json:"ingredient_id" gorm:"not null"`       // FK a product_product
	Quantity          float64 `json:"quantity" gorm:"type:decimal(10,4);not null"`
	UnitID            uint    `json:"unit_id" gorm:"not null"`
	WastePercentage   float64 `json:"waste_percentage" gorm:"type:decimal(5,2);default:0;not null"`
	Notes             string  `json:"notes" gorm:"type:text"`

	// Relaciones
	ProductTemplate *ProductTemplate `json:"product_template,omitempty" gorm:"foreignKey:ProductTemplateID"`
	Ingredient      *ProductProduct  `json:"ingredient,omitempty" gorm:"foreignKey:IngredientID"`
	Unit            *Unit            `json:"unit,omitempty" gorm:"foreignKey:UnitID"`
}

func (Recipe) TableName() string {
	return "recipes"
}
