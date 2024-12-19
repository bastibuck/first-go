package entities

import (
	"gorm.io/gorm"
)

type Users struct {
	gorm.Model
	Email         string
	Password_Hash string
}
