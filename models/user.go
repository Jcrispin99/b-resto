package models

import "gorm.io/gorm"

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
	GuestRole Role = "guest"
)

type User struct {
	gorm.Model

	// Relación con Company
	CompanyID *uint    `json:"company_id" gorm:"index"`
	Company   *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`

	// Información del usuario
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" gorm:"unique;not null" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     Role   `json:"role" binding:"required,oneof=admin user guest"`
	IsActive bool   `json:"is_active" gorm:"default:true"`
}
