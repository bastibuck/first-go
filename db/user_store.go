package db

import (
	"context"
	"first-go/db/entities"
	userTypes "first-go/types/user"

	"gorm.io/gorm"
)

type UserStore interface {
	Create(ctx context.Context, user *userTypes.User) error
	GetByEmail(ctx context.Context, email string) (*userTypes.User, error)
}

type DatabaseUserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *DatabaseUserStore {
	return &DatabaseUserStore{
		db,
	}
}

func (store *DatabaseUserStore) Create(ctx context.Context, user *userTypes.User) error {
	var newUser = entities.Users{
		Email:         user.Email,
		Password_Hash: user.PasswordHash,
	}

	result := store.db.WithContext(ctx).Create(&newUser)

	return result.Error
}

func (store *DatabaseUserStore) GetByEmail(ctx context.Context, email string) (*userTypes.User, error) {
	var user entities.Users

	result := store.db.WithContext(ctx).Where("email = ?", email).First(&user)

	if result.Error != nil {
		return nil, result.Error
	}

	userResponse := userTypes.User{
		ID:           user.ID,
		Email:        user.Email,
		PasswordHash: user.Password_Hash,
	}

	return &userResponse, nil

}
