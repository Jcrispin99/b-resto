package models

import "time"

// CashMovement - Movimiento de efectivo durante una sesión
type CashMovement struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	POSSessionID  uint      `json:"pos_session_id" gorm:"not null;column:cash_register_id"`
	Type          string    `json:"type" gorm:"size:50;not null"`     // "in" (entrada), "out" (salida)
	Concept       string    `json:"concept" gorm:"size:255;not null"` // "Venta", "Gasto", "Retiro"
	Amount        float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	PaymentMethod string    `json:"payment_method" gorm:"size:50"` // "cash", "card", "transfer"
	Reference     string    `json:"reference" gorm:"size:255"`     // Referencia/Voucher
	UserID        uint      `json:"user_id" gorm:"not null"`
	Notes         string    `json:"notes" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`

	// Relaciones
	POSSession *POSSession `json:"pos_session,omitempty" gorm:"foreignKey:POSSessionID"`
	User       *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

func (CashMovement) TableName() string {
	return "cash_movements"
}
