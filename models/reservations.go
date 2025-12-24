package models

import (
	"time"

	"gorm.io/gorm"
)

// Reservation - Reservas de mesas
type Reservation struct {
	gorm.Model
	ReservationNumber string     `json:"reservation_number" gorm:"size:100;not null;uniqueIndex"`
	CompanyID         uint       `json:"company_id" gorm:"not null"`
	PartnerID         uint       `json:"partner_id" gorm:"not null"` // Cliente (de tabla partners)
	TableID           *uint      `json:"table_id"`                   // Nullable - se asigna después
	ReservationDate   time.Time  `json:"reservation_date" gorm:"type:date;not null"`
	ReservationTime   time.Time  `json:"reservation_time" gorm:"type:time;not null"`
	GuestsCount       int        `json:"guests_count" gorm:"not null" binding:"required,min=1"`
	Status            string     `json:"status" gorm:"size:50;default:'pending';not null"` // pending, confirmed, seated, cancelled, no_show
	SpecialRequests   string     `json:"special_requests" gorm:"type:text"`
	ConfirmedAt       *time.Time `json:"confirmed_at"`
	SeatedAt          *time.Time `json:"seated_at"`

	// Relaciones
	Company *Company `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Partner *Partner `json:"partner,omitempty" gorm:"foreignKey:PartnerID"` // Cliente
	Table   *Table   `json:"table,omitempty" gorm:"foreignKey:TableID"`
}

func (Reservation) TableName() string {
	return "reservations"
}
