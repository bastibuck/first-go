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
	GetAll(ctx context.Context) ([]eventTypes.EventListResponse, error)
	GetById(ctx context.Context, id uint) (*eventTypes.EventResponse, error)
	AddEvent(ctx context.Context, event *eventTypes.EventUpsertPayload, userID uint) error
	UpdateEvent(ctx context.Context, id uint, event *eventTypes.EventUpsertPayload) error
	DeleteById(ctx context.Context, id uint) error
	SignUp(ctx context.Context, eventId uint, signUp *entities.EventSignUps) error
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

func (store *DatabaseEventStore) GetAll(ctx context.Context) ([]eventTypes.EventListResponse, error) {
	var events []eventTypes.EventListResponse

	result := store.db.WithContext(ctx).Table("events").Find(&events).Order("date desc")
	if result.Error != nil {
		return nil, result.Error
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

	event := entities.Events{Model: gorm.Model{ID: id}}
	result := store.db.WithContext(ctx).Model(&event).Select("events.id, events.name, events.date, events.description, events.pax, users.id as user_id, users.email as user_email").Joins("join users on events.user_id = users.id").Scan(&eventsResult)
	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var signupsResult []struct {
		Email string `json:"email"`
	}

	result = store.db.WithContext(ctx).Model(&entities.EventSignUps{}).Where("event_sign_ups.event_id = ?", id).Select("users.email").Joins("join users on event_sign_ups.user_id = users.id").Scan(&signupsResult)
	if result.Error != nil {
		return nil, result.Error
	}

	signupsResponse := []string{}

	for _, signup := range signupsResult {
		signupsResponse = append(signupsResponse, signup.Email)
	}

	eventResponse := eventTypes.EventResponse{
		ID:          eventsResult.ID,
		Name:        eventsResult.Name,
		Description: eventsResult.Description,
		Date:        eventsResult.Date,
		Pax:         eventsResult.Pax,
		User: userTypes.User{
			ID:    eventsResult.UserID,
			Email: eventsResult.UserEmail,
		},
		SignUps: signupsResponse,
	}

	return &eventResponse, nil
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
	var existingEvent entities.Events

	// Fetch the existing event
	result := store.db.WithContext(ctx).First(&existingEvent, id)
	if result.Error != nil {
		return result.Error
	}

	// Update the fields
	existingEvent.Name = event.Name
	existingEvent.Date = event.Date
	existingEvent.Description = event.Description
	existingEvent.Pax = event.Pax

	// Save the updated event
	result = store.db.WithContext(ctx).Save(&existingEvent)

	return result.Error
}

func (store *DatabaseEventStore) DeleteById(ctx context.Context, id uint) error {
	var event entities.Events

	result := store.db.WithContext(ctx).Where("id = ?", id).Unscoped().Delete(&event)

	return result.Error
}

func (store *DatabaseEventStore) SignUp(ctx context.Context, eventId uint, signUp *entities.EventSignUps) error {
	event, err := store.GetById(ctx, eventId)
	if err != nil {
		return err
	}

	if event.Pax <= 0 {
		return fmt.Errorf("event-fully-booked")
	}

	err = store.UpdateEvent(ctx, eventId, &eventTypes.EventUpsertPayload{
		Name:        event.Name,
		Date:        event.Date,
		Description: event.Description,
		Pax:         event.Pax - 1,
	})

	if err != nil {
		return err
	}

	result := store.db.WithContext(ctx).Create(&signUp)

	return result.Error
}
