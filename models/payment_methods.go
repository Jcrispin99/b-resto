package models

import "gorm.io/gorm"

// PaymentMethod representa un método de pago
type PaymentMethod struct {
	gorm.Model
	Code     string `json:"code" gorm:"size:20;not null;unique" binding:"required,min=2,max=20"`
	Name     string `json:"name" gorm:"size:50;not null" binding:"required,min=3,max=50"`
	Type     string `json:"type" gorm:"size:255;default:'cash';not null" binding:"required,oneof=cash transfer card wallet"`
	IsActive bool   `json:"is_active" gorm:"default:true;not null"`
}

// TableName especifica el nombre de la tabla
func (PaymentMethod) TableName() string {
	return "payment_methods"
}
