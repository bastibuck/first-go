package entities

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex"`
	Password_Hash string
}
