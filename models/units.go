package models

import "gorm.io/gorm"

type Unit struct {
	gorm.Model
	Name         string `json:"name" gorm:"not null" binding:"required,min=2,max=100"`
	Abbreviation string `json:"abbreviation" gorm:"not null" binding:"required,min=1,max=10"`
	Type         string `json:"type" gorm:"not null" binding:"required,oneof=weight volume unit length area"`
	IsActive     bool   `json:"is_active" gorm:"default:true"`
}
