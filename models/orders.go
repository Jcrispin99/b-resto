package models

import (
	"time"

	"gorm.io/gorm"
)

// Order - Órdenes de venta del POS
type Order struct {
	gorm.Model
	JournalID   uint      `json:"journal_id" gorm:"not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	TableID     *uint     `json:"table_id"`                                      // Nullable - null si es para llevar
	Name        string    `json:"name" gorm:"size:100;not null"`                 // SO/2024/0001
	State       string    `json:"state" gorm:"size:50;default:'draft';not null"` // draft, confirmed, done, cancelled
	OrderDate   time.Time `json:"order_date" gorm:"type:date;not null"`
	TotalAmount float64   `json:"total_amount" gorm:"type:decimal(10,2);default:0;not null"`
	Note        string    `json:"note" gorm:"type:text"`

	// Relaciones
	Journal  *Journal        `json:"journal,omitempty" gorm:"foreignKey:JournalID"`
	User     *User           `json:"user,omitempty" gorm:"foreignKey:UserID"`
	Items    []OrderItem     `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	Payments []OrderPayment  `json:"payments,omitempty" gorm:"foreignKey:OrderID"`
	Tickets  []KitchenTicket `json:"tickets,omitempty" gorm:"foreignKey:OrderID"`
}

func (Order) TableName() string {
	return "orders"
}
