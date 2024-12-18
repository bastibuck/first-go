package types

import "time"

type EventResponse struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

type EventPayloadUpsert struct {
	Name        string    `json:"name"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}
