package models

import "gorm.io/gorm"

// Tax representa un impuesto
type Tax struct {
	gorm.Model
	Name             string  `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	TaxType          string  `json:"tax_type" gorm:"size:255;not null" binding:"required,max=255"`
	RatePercent      float64 `json:"rate_percent" gorm:"type:decimal(5,2);default:0;not null" binding:"omitempty,gte=0,lte=100"`
	IsPriceInclusive bool    `json:"is_price_inclusive" gorm:"default:false;not null"`
	IsActive         bool    `json:"is_active" gorm:"default:true;not null"`
}

// TableName especifica el nombre de la tabla
func (Tax) TableName() string {
	return "taxes"
}
