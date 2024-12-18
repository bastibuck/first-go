package types

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	PasswordHash string `json:"-"` // - omits the output when reading
}

type NewUserPayload struct {
	Email    string
	Password string
}

func NewUser(createUser NewUserPayload) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(createUser.Password), 14)

	if err != nil {
		return nil, err
	}

	return &User{
		Email:        createUser.Email,
		PasswordHash: string(hashedPassword),
	}, nil
}

func ValidatePassword(hashed string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(password)) == nil
}

func CreateToken(user User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 4).Unix(), // TODO! does this work without manually checking?
	}, nil)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
