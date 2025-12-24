package models

import "gorm.io/gorm"

// Partner - Proveedores/Clientes/Socios comerciales
type Partner struct {
	gorm.Model
	Code             string `json:"code" gorm:"size:50;not null;uniqueIndex"`
	PartnerType      string `json:"partner_type" gorm:"size:50;default:'company';not null"` // company, person
	Name             string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	TradeName        string `json:"trade_name" gorm:"size:255"`
	TaxID            string `json:"tax_id" gorm:"size:50;not null"`
	Email            string `json:"email" gorm:"size:255;not null"`
	Phone            string `json:"phone" gorm:"size:50;not null"`
	Address          string `json:"address" gorm:"size:500"`
	UbigeoCode       string `json:"ubigeo_code" gorm:"size:10"` // Código geográfico
	IsCustomer       bool   `json:"is_customer" gorm:"default:false;not null"`
	IsSupplier       bool   `json:"is_supplier" gorm:"default:false;not null"`
	PaymentTermsDays int    `json:"payment_terms_days" gorm:"default:0;not null"` // Días de crédito
	Notes            string `json:"notes" gorm:"type:text"`
	IsActive         bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	PurchaseOrders []PurchaseOrder `json:"purchase_orders,omitempty" gorm:"foreignKey:PartnerID"`
}

func (Partner) TableName() string {
	return "partners"
}
