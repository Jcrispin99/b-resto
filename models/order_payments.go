package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderPayment struct {
	gorm.Model
	OrderID         uint      `json:"order_id" gorm:"not null"`
	PaymentMethodID uint      `json:"payment_method_id" gorm:"not null"`
	JournalID       uint      `json:"journal_id" gorm:"not null"` // Journal de caja
	Amount          float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentDate     time.Time `json:"payment_date" gorm:"type:date;not null"`

	// Relaciones
	Order         *Order         `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	PaymentMethod *PaymentMethod `json:"payment_method,omitempty" gorm:"foreignKey:PaymentMethodID"`
	Journal       *Journal       `json:"journal,omitempty" gorm:"foreignKey:JournalID"`
}

func (OrderPayment) TableName() string {
	return "order_payments"
}
