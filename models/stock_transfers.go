package models

import (
	"time"

	"gorm.io/gorm"
)

// StockTransfer - Transferencias de stock entre almacenes
type StockTransfer struct {
	gorm.Model
	TransferNumber  string    `json:"transfer_number" gorm:"size:100;not null"`
	JournalID       *uint     `json:"journal_id"`
	FromWarehouseID uint      `json:"from_warehouse_id" gorm:"not null"`
	ToWarehouseID   uint      `json:"to_warehouse_id" gorm:"not null"`
	TransferDate    time.Time `json:"transfer_date" gorm:"type:date;not null"`
	Status          string    `json:"status" gorm:"size:50;default:'pending';not null"` // pending, in_transit, received
	Notes           string    `json:"notes" gorm:"type:text"`
	CreatedBy       *uint     `json:"created_by"`
	ReceivedBy      *uint     `json:"received_by"`

	// Relaciones
	Journal        *Journal            `json:"journal,omitempty" gorm:"foreignKey:JournalID"`
	FromWarehouse  *Warehouse          `json:"from_warehouse,omitempty" gorm:"foreignKey:FromWarehouseID"`
	ToWarehouse    *Warehouse          `json:"to_warehouse,omitempty" gorm:"foreignKey:ToWarehouseID"`
	CreatedByUser  *User               `json:"created_by_user,omitempty" gorm:"foreignKey:CreatedBy"`
	ReceivedByUser *User               `json:"received_by_user,omitempty" gorm:"foreignKey:ReceivedBy"`
	Items          []StockTransferItem `json:"items,omitempty" gorm:"foreignKey:StockTransferID"`
}

func (StockTransfer) TableName() string {
	return "stock_transfers"
}
