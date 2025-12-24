package models

import (
	"time"

	"gorm.io/gorm"
)

// PurchaseOrder - Órdenes de compra a proveedores
type PurchaseOrder struct {
	gorm.Model
	OrderNumber          string     `json:"order_number" gorm:"size:100;not null"`
	JournalID            *uint      `json:"journal_id"`
	CompanyID            uint       `json:"company_id" gorm:"not null"`
	WarehouseID          uint       `json:"warehouse_id" gorm:"not null"`
	PartnerID            uint       `json:"partner_id" gorm:"not null"`
	OrderDate            time.Time  `json:"order_date" gorm:"type:date;not null"`
	ExpectedDeliveryDate *time.Time `json:"expected_delivery_date" gorm:"type:date"`
	ReceivedDate         *time.Time `json:"received_date" gorm:"type:date"`
	PaidDate             *time.Time `json:"paid_date" gorm:"type:date"`
	Status               string     `json:"status" gorm:"size:50;default:'quote_request';not null"` // quote_request, confirmed, received, paid
	Subtotal             float64    `json:"subtotal" gorm:"type:decimal(10,2);not null"`
	Tax                  float64    `json:"tax" gorm:"type:decimal(10,2);not null"`
	Total                float64    `json:"total" gorm:"type:decimal(10,2);not null"`
	Notes                string     `json:"notes" gorm:"type:text"`
	CreatedBy            *uint      `json:"created_by"`
	ApprovedBy           *uint      `json:"approved_by"`

	// Relaciones
	Journal        *Journal            `json:"journal,omitempty" gorm:"foreignKey:JournalID"`
	Company        *Company            `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Warehouse      *Warehouse          `json:"warehouse,omitempty" gorm:"foreignKey:WarehouseID"`
	Partner        *Partner            `json:"partner,omitempty" gorm:"foreignKey:PartnerID"`
	CreatedByUser  *User               `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBy"`
	ApprovedByUser *User               `json:"approved_by_user,omitempty" gorm:"foreignKey:ApprovedBy"`
	Items          []PurchaseOrderItem `json:"items,omitempty" gorm:"foreignKey:PurchaseOrderID"`
}

func (PurchaseOrder) TableName() string {
	return "purchase_orders"
}
