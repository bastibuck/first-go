package db

import (
	"context"
	"database/sql"

	eventTypes "first-go/types/event"
)

type EventStore interface {
	GetAll(ctx context.Context) ([]eventTypes.EventResponse, error)
	GetById(ctx context.Context, id int) (*eventTypes.EventResponse, error)
	AddEvent(ctx context.Context, event *eventTypes.EventPayloadUpsert, userID int) error
	UpdateEvent(ctx context.Context, id int, event *eventTypes.EventPayloadUpsert) error
	DeleteById(ctx context.Context, id int) error
}

type DatabaseEventStore struct {
	db *sql.DB
}

// -----------------
// Constructor for EventStore
// -----------------

func NewEventStore(db *sql.DB) *DatabaseEventStore {
	return &DatabaseEventStore{
		db,
	}
}

// -----------------
// Functions to interact with event in the database
// -----------------

func (store *DatabaseEventStore) GetAll(ctx context.Context) ([]eventTypes.EventResponse, error) {
	query := `
		SELECT e.id, e.name, e.date, e.description, u.id, u.email
		FROM events e
		JOIN users u ON e.user_id = u.id
		ORDER BY e.date DESC
	`

	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []eventTypes.EventResponse{}

	for rows.Next() {
		var event eventTypes.EventResponse

		err := rows.Scan(&event.ID, &event.Name, &event.Date, &event.Description, &event.User.ID, &event.User.Email)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (store *DatabaseEventStore) GetById(ctx context.Context, id int) (*eventTypes.EventResponse, error) {
	query := `
		SELECT e.id, e.name, e.date, e.description, u.id, u.email
		FROM events e
		JOIN users u ON e.user_id = u.id
		WHERE e.id = ?
	`

	row := store.db.QueryRowContext(ctx, query, id)

	var event eventTypes.EventResponse

	err := row.Scan(&event.ID, &event.Name, &event.Date, &event.Description, &event.User.ID, &event.User.Email)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (store *DatabaseEventStore) AddEvent(ctx context.Context, event *eventTypes.EventPayloadUpsert, userID int) error {
	insertEventSQL := `INSERT INTO events(name, date, description, user_id) VALUES (?, ?, ?, ?)`

	statement, err := store.db.Prepare(insertEventSQL)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, event.Name, event.Date, event.Description, userID)

	return err
}

func (store *DatabaseEventStore) UpdateEvent(ctx context.Context, id int, event *eventTypes.EventPayloadUpsert) error {
	updateEventSQL := `UPDATE events SET name = ?, date = ?, description = ? WHERE id = ?`

	statement, err := store.db.Prepare(updateEventSQL)
	if err != nil {
		return err
	}

	result, err := statement.ExecContext(ctx, event.Name, event.Date, event.Description, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (store *DatabaseEventStore) DeleteById(ctx context.Context, id int) error {
	deleteEventSQL := `DELETE FROM events WHERE id = ?`

	statement, err := store.db.Prepare(deleteEventSQL)
	if err != nil {
		return err
	}

	result, err := statement.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
