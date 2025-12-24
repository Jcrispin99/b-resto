package models

import (
	"time"

	"gorm.io/gorm"
)

type KitchenTicket struct {
	gorm.Model
	OrderID          uint      `json:"order_id" gorm:"not null"`
	KitchenStationID uint      `json:"kitchen_station_id" gorm:"not null"`
	TicketNumber     string    `json:"ticket_number" gorm:"size:50;not null"`
	State            string    `json:"state" gorm:"size:50;default:'pending'"` // pending, preparing, ready, delivered
	CreatedDate      time.Time `json:"created_date"`

	// Relaciones
	Order          *Order              `json:"order,omitempty" gorm:"foreignKey:OrderID"`
	KitchenStation *KitchenStation     `json:"kitchen_station,omitempty" gorm:"foreignKey:KitchenStationID"`
	Items          []KitchenTicketItem `json:"items,omitempty" gorm:"foreignKey:KitchenTicketID"`
}

func (KitchenTicket) TableName() string {
	return "kitchen_tickets"
}
