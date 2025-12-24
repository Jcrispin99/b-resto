package models

import "gorm.io/gorm"

// POS - Punto de Venta físico (terminal, computadora, tablet)
type POS struct {
	gorm.Model
	CompanyID        uint   `json:"company_id" gorm:"not null"`
	Code             string `json:"code" gorm:"size:50;not null;uniqueIndex"`
	Name             string `json:"name" gorm:"size:255;not null" binding:"required,min=3,max=255"`
	IPAddress        string `json:"ip_address" gorm:"size:50"`
	PrinterIP        string `json:"printer_ip" gorm:"size:50"`
	DefaultJournalID *uint  `json:"default_journal_id"`
	IsActive         bool   `json:"is_active" gorm:"default:true;not null"`

	// Relaciones
	Company        *Company     `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	DefaultJournal *Journal     `json:"default_journal,omitempty" gorm:"foreignKey:DefaultJournalID"`
	Sessions       []POSSession `json:"sessions,omitempty" gorm:"foreignKey:POSID"`
}

func (POS) TableName() string {
	return "pos_terminals"
}
