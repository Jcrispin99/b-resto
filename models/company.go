package models

import "gorm.io/gorm"

// Company representa una compañía o sucursal
type Company struct {
	gorm.Model

	// Relación jerárquica (auto-referencia para sucursales)
	ParentID *uint     `json:"parent_id" gorm:"index"`
	Parent   *Company  `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Branches []Company `json:"branches,omitempty" gorm:"foreignKey:ParentID"`

	// Información básica
	Name         string `json:"name" gorm:"size:200;not null" binding:"required,min=3,max=200"`
	BusinessName string `json:"business_name" gorm:"size:200;not null" binding:"required,min=3,max=200"`
	Logo         string `json:"logo" gorm:"size:255" binding:"omitempty,max=255"`
	IsActive     bool   `json:"is_active" gorm:"default:true;not null"`

	// Contacto
	Phone   string `json:"phone" gorm:"size:20" binding:"omitempty,min=7,max=20"`
	Email   string `json:"email" gorm:"size:150" binding:"omitempty,email,max=150"`
	Website string `json:"website" gorm:"size:255" binding:"omitempty,url,max=255"`

	// Ubicación
	Address    string `json:"address" gorm:"type:text" binding:"omitempty,max=500"`
	UbigeoCode string `json:"ubigeo_code" gorm:"size:6" binding:"omitempty,len=6"`
}

// TableName especifica el nombre de la tabla
func (Company) TableName() string {
	return "companies"
}

// IsHeadquarters verifica si es la casa matriz
func (c *Company) IsHeadquarters() bool {
	return c.ParentID == nil
}

// IsBranch verifica si es una sucursal
func (c *Company) IsBranch() bool {
	return c.ParentID != nil
}
