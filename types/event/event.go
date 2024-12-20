package types

import (
	userTypes "first-go/types/user"
	"time"
)

type EventResponse struct {
	ID          uint           `json:"id"`
	Name        string         `json:"name"`
	Date        time.Time      `json:"date"`
	Description string         `json:"description"`
	Pax         int            `json:"pax"`
	User        userTypes.User `json:"user"`
}

type EventUpsertPayload struct {
	Name        string    `json:"name" validate:"required,min=3"`
	Date        time.Time `json:"date" validate:"required"` // TODO? how to validate date?
	Description string    `json:"description" validate:"required,max=255"`
	Pax         int       `json:"pax" validate:"required,min=1"`
}

type EventSignUpPayload struct {
	Email string `json:"email" validate:"required,email"`
	Name  string `json:"name" validate:"required,min=3"`
}
