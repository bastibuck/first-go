package entities

import (
	"time"

	"gorm.io/gorm"
)

type Events struct {
	gorm.Model
	Name        string `gorm:"uniqueIndex"`
	Description string
	Date        time.Time
	UserID      uint
	Pax         int
}

type EventSignUps struct {
	gorm.Model
	UserID  uint `gorm:"uniqueIndex:idx_event_to_user"`
	EventID uint `gorm:"uniqueIndex:idx_event_to_user"`
}
