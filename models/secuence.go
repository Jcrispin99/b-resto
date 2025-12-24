package models

import "gorm.io/gorm"

type Sequence struct {
	gorm.Model
	JournalID  uint   `json:"journal_id" gorm:"not null"`
	Name       string `json:"name" gorm:"size:255;not null"`
	Padding    int    `json:"padding" gorm:"default:5"`
	NextNumber int    `json:"next_number" gorm:"default:1"`
	IsActive   bool   `json:"is_active" gorm:"default:true"`

	// Relaciones
	Journal *Journal `json:"journal,omitempty" gorm:"foreignKey:JournalID"`
}

func (Sequence) TableName() string {
	return "sequences"
}
