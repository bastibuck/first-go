package db

import (
	"context"
	"database/sql"
	"first-go/types"
)

type UserStore interface {
	Create(ctx context.Context, user *types.User) error
	GetByEmail(ctx context.Context, email string) (*types.User, error)
}

type DatabaseUserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *DatabaseUserStore {
	return &DatabaseUserStore{
		db,
	}
}

func (store *DatabaseUserStore) Create(ctx context.Context, user *types.User) error {
	insertUserSQL := `
		INSERT INTO users (email, password_hash)
		VALUES (?,?)
	`

	statement, err := store.db.Prepare(insertUserSQL)
	if err != nil {
		return err
	}

	_, err = statement.ExecContext(ctx, user.Email, user.PasswordHash)

	return err
}

func (u *DatabaseUserStore) GetByEmail(ctx context.Context, email string) (*types.User, error) {
	query := `
		SELECT id, email, password_hash FROM users WHERE email = ?
	`

	var user types.User

	err := u.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
