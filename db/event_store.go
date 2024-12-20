package db

import (
	"context"
	"fmt"
	"time"

	"first-go/db/entities"
	eventTypes "first-go/types/event"
	userTypes "first-go/types/user"

	"gorm.io/gorm"
)

type EventStore interface {
	GetAll(ctx context.Context) ([]eventTypes.EventResponse, error)
	GetById(ctx context.Context, id uint) (*eventTypes.EventResponse, error)
	AddEvent(ctx context.Context, event *eventTypes.EventUpsertPayload, userID uint) error
	UpdateEvent(ctx context.Context, id uint, event *eventTypes.EventUpsertPayload) error
	DeleteById(ctx context.Context, id uint) error
	SignUp(ctx context.Context, eventId uint, signUp *eventTypes.EventSignUpPayload) error
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
		Pax         int       `json:"pax"`
		UserID      uint      `json:"user_id"`
		UserEmail   string    `json:"user_email"`
	}

	result := store.db.WithContext(ctx).Raw(`
		SELECT e.id, e.name, e.date, e.description, e.pax, u.id as user_id, u.email as user_email
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
			Pax:         eventsResult[event].Pax,
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
		Pax         int       `json:"pax"`
		UserID      uint      `json:"user_id"`
		UserEmail   string    `json:"user_email"`
	}

	result := store.db.WithContext(ctx).Raw(`
		SELECT e.id, e.name, e.date, e.description, e.pax, u.id as user_id, u.email as user_email
		FROM events e
		JOIN users u ON e.user_id = u.id
		ORDER BY e.date DESC
		LIMIT 1
	`).Scan(&eventsResult)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	event := eventTypes.EventResponse{
		ID:          eventsResult.ID,
		Name:        eventsResult.Name,
		Description: eventsResult.Description,
		Date:        eventsResult.Date,
		Pax:         eventsResult.Pax,
		User: userTypes.User{
			ID:    eventsResult.UserID,
			Email: eventsResult.UserEmail,
		},
	}

	return &event, nil
}

func (store *DatabaseEventStore) AddEvent(ctx context.Context, event *eventTypes.EventUpsertPayload, userID uint) error {
	var newEvent = entities.Events{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
		UserID:      userID,
		Pax:         event.Pax,
	}

	result := store.db.WithContext(ctx).Create(&newEvent)

	return result.Error
}

func (store *DatabaseEventStore) UpdateEvent(ctx context.Context, id uint, event *eventTypes.EventUpsertPayload) error {
	var updatedEvent = entities.Events{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
		Pax:         event.Pax,
	}

	result := store.db.WithContext(ctx).Model(&updatedEvent).Where("id = ?", id).Updates(&updatedEvent)

	return result.Error
}

func (store *DatabaseEventStore) DeleteById(ctx context.Context, id uint) error {
	var event entities.Events

	result := store.db.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&event)

	return result.Error
}

func (store *DatabaseEventStore) SignUp(ctx context.Context, eventId uint, signUp *eventTypes.EventSignUpPayload) error {
	var newSignUp = entities.EventSignUps{
		Email:   signUp.Email,
		Name:    signUp.Name,
		EventID: eventId,
	}

	event, err := store.GetById(ctx, eventId)
	if err != nil {
		return err
	}

	if event.Pax <= 0 {
		return fmt.Errorf("event-fully-booked")
	}

	event.Pax--

	err = store.UpdateEvent(ctx, eventId, &eventTypes.EventUpsertPayload{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
		Pax:         event.Pax,
	})

	if err != nil {
		return err
	}

	result := store.db.WithContext(ctx).Create(&newSignUp)

	return result.Error
}
