package db

import (
	"context"
	"database/sql"

	"first-go/types"
)

type EventStore interface {
	GetAll(ctx context.Context) ([]types.EventResponse, error)
	GetById(ctx context.Context, id int) (*types.EventResponse, error)
	AddEvent(ctx context.Context, event *types.CreateEvent) error
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

func (store *DatabaseEventStore) GetAll(ctx context.Context) ([]types.EventResponse, error) {
	query := `SELECT id, name, date, description FROM events ORDER BY date DESC`

	rows, err := store.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	events := []types.EventResponse{}

	for rows.Next() {
		var event types.EventResponse

		err := rows.Scan(&event.ID, &event.Name, &event.Date, &event.Description)
		if err != nil {
			return nil, err
		}

		events = append(events, event)
	}

	return events, nil
}

func (store *DatabaseEventStore) GetById(ctx context.Context, id int) (*types.EventResponse, error) {
	query := `SELECT id, name, date, description FROM events WHERE id = ?`

	row := store.db.QueryRowContext(ctx, query, id)

	var event types.EventResponse

	err := row.Scan(&event.ID, &event.Name, &event.Date, &event.Description)
	if err != nil {
		return nil, err
	}

	return &event, nil
}

func (store *DatabaseEventStore) AddEvent(ctx context.Context, event *types.CreateEvent) error {
	insertEventSQL := `INSERT INTO events(name, date, description) VALUES (?, ?, ?)`

	statement, err := store.db.Prepare(insertEventSQL)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, event.Name, event.Date, event.Description)

	return err
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
