package models

import "gorm.io/gorm"

type Journal struct {
	gorm.Model
	CompanyID uint   `json:"company_id" gorm:"not null"`
	Code      string `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Name      string `json:"name" gorm:"size:255;not null"`
	Type      string `json:"type" gorm:"size:50;not null"` // sale, purchase, cash, bank
	IsActive  bool   `json:"is_active" gorm:"default:true"`

	// Relaciones
	Company   *Company   `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Sequences []Sequence `json:"sequences,omitempty" gorm:"foreignKey:JournalID"`
}

func (Journal) TableName() string {
	return "journals"
}
