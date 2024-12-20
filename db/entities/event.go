package entities

import (
	"time"

	"gorm.io/gorm"
)

type Events struct {
	gorm.Model
	Name        string
	Description string
	Date        time.Time
	UserID      uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Pax         int
}

type EventSignUps struct {
	gorm.Model
	Email   string
	Name    string
	EventID uint
}
