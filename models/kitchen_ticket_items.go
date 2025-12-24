package models

import "gorm.io/gorm"

type KitchenTicketItem struct {
	gorm.Model
	KitchenTicketID uint    `json:"kitchen_ticket_id" gorm:"not null"`
	OrderItemID     uint    `json:"order_item_id" gorm:"not null"`
	Quantity        float64 `json:"quantity" gorm:"type:decimal(10,2);not null"`

	// Relaciones
	KitchenTicket *KitchenTicket `json:"kitchen_ticket,omitempty" gorm:"foreignKey:KitchenTicketID"`
	OrderItem     *OrderItem     `json:"order_item,omitempty" gorm:"foreignKey:OrderItemID"`
}

func (KitchenTicketItem) TableName() string {
	return "kitchen_ticket_items"
}
