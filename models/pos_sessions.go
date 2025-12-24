package models

import (
	"time"

	"gorm.io/gorm"
)

// POSSession - Sesión de caja (turno de apertura a cierre)
type POSSession struct {
	gorm.Model
	POSID           uint       `json:"pos_id" gorm:"not null"`
	OpeningBalance  float64    `json:"opening_balance" gorm:"type:decimal(10,2);not null"`
	ClosingBalance  *float64   `json:"closing_balance" gorm:"type:decimal(10,2)"`
	ExpectedBalance *float64   `json:"expected_balance" gorm:"type:decimal(10,2)"`
	Difference      *float64   `json:"difference" gorm:"type:decimal(10,2)"` // Faltante/Sobrante
	OpenedBy        uint       `json:"opened_by" gorm:"not null"`
	ClosedBy        *uint      `json:"closed_by"`
	OpenedAt        time.Time  `json:"opened_at" gorm:"not null"`
	ClosedAt        *time.Time `json:"closed_at"`
	Status          string     `json:"status" gorm:"size:50;default:'open';not null"` // open, closed
	Notes           string     `json:"notes" gorm:"type:text"`

	// Relaciones
	POS          *POS           `json:"pos,omitempty" gorm:"foreignKey:POSID"`
	OpenedByUser *User          `json:"opened_by_user,omitempty" gorm:"foreignKey:OpenedBy"`
	ClosedByUser *User          `json:"closed_by_user,omitempty" gorm:"foreignKey:ClosedBy"`
	Movements    []CashMovement `json:"movements,omitempty" gorm:"foreignKey:POSSessionID"`
}

func (POSSession) TableName() string {
	return "cash_registers"
}
