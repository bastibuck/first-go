package types

import "time"

type EventResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	UserID      int       `json:"user_id"` // TODO: map to user at some point
}

type EventPayloadUpsert struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}
