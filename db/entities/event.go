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
	Pax         int
}

type EventSignUps struct {
	gorm.Model
	UserID  uint
	EventID uint
}
