package db

import (
	"context"
	"time"

	"first-go/db/entities"
	eventTypes "first-go/types/event"
	userTypes "first-go/types/user"

	"gorm.io/gorm"
)

type EventStore interface {
	GetAll(ctx context.Context) ([]eventTypes.EventResponse, error)
	GetById(ctx context.Context, id uint) (*eventTypes.EventResponse, error)
	AddEvent(ctx context.Context, event *eventTypes.EventPayloadUpsert, userID uint) error
	UpdateEvent(ctx context.Context, id uint, event *eventTypes.EventPayloadUpsert) error
	DeleteById(ctx context.Context, id uint) error
}

type DatabaseEventStore struct {
	db *gorm.DB
}

// -----------------
// Constructor for EventStore
// -----------------

func NewEventStore(db *gorm.DB) *DatabaseEventStore {
	return &DatabaseEventStore{
		db,
	}
}

// -----------------
// Functions to interact with event in the database
// -----------------

func (store *DatabaseEventStore) GetAll(ctx context.Context) ([]eventTypes.EventResponse, error) {
	var eventsResult []struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		Date        time.Time `json:"date"`
		Description string    `json:"description"`
		UserID      uint      `json:"user_id"`
		UserEmail   string    `json:"user_email"`
	}

	result := store.db.WithContext(ctx).Raw(`
		SELECT e.id, e.name, e.date, e.description, u.id as user_id, u.email as user_email
		FROM events e
		JOIN users u ON e.user_id = u.id
		ORDER BY e.date DESC
	`).Scan(&eventsResult)

	if result.Error != nil {
		return nil, result.Error
	}

	events := make([]eventTypes.EventResponse, len(eventsResult))

	for event := range events {
		events[event] = eventTypes.EventResponse{
			ID:          eventsResult[event].ID,
			Name:        eventsResult[event].Name,
			Date:        eventsResult[event].Date,
			Description: eventsResult[event].Description,
			User: userTypes.User{
				ID:    eventsResult[event].UserID,
				Email: eventsResult[event].UserEmail,
			},
		}
	}

	return events, nil
}

func (store *DatabaseEventStore) GetById(ctx context.Context, id uint) (*eventTypes.EventResponse, error) {
	var eventsResult struct {
		ID          uint      `json:"id"`
		Name        string    `json:"name"`
		Date        time.Time `json:"date"`
		Description string    `json:"description"`
		UserID      uint      `json:"user_id"`
		UserEmail   string    `json:"user_email"`
	}

	result := store.db.WithContext(ctx).Raw(`
		SELECT e.id, e.name, e.date, e.description, u.id as user_id, u.email as user_email
		FROM events e
		JOIN users u ON e.user_id = u.id
		ORDER BY e.date DESC
		LIMIT 1
	`).Scan(&eventsResult)

	if result.Error != nil {
		return nil, result.Error
	}

	event := eventTypes.EventResponse{
		ID:          eventsResult.ID,
		Name:        eventsResult.Name,
		Description: eventsResult.Description,
		Date:        eventsResult.Date,
		User: userTypes.User{
			ID:    eventsResult.UserID,
			Email: eventsResult.UserEmail,
		},
	}

	return &event, nil
}

func (store *DatabaseEventStore) AddEvent(ctx context.Context, event *eventTypes.EventPayloadUpsert, userID uint) error {
	var newEvent = entities.Events{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
		UserID:      userID,
	}

	result := store.db.WithContext(ctx).Create(&newEvent)

	return result.Error
}

func (store *DatabaseEventStore) UpdateEvent(ctx context.Context, id uint, event *eventTypes.EventPayloadUpsert) error {
	var updatedEvent = entities.Events{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
	}

	result := store.db.WithContext(ctx).Model(&updatedEvent).Where("id = ?", id).Updates(&updatedEvent)

	return result.Error
}

func (store *DatabaseEventStore) DeleteById(ctx context.Context, id uint) error {
	var event entities.Events

	result := store.db.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&event)

	return result.Error
}
