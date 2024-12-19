package middleware

import (
	"context"
	userTypes "first-go/types/user"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type key string

const (
	userKey key = "user"
)

func WithUser(ctx context.Context, user *userTypes.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// get User from context
func User(ctx context.Context) *userTypes.User {
	val := ctx.Value(userKey)
	user, ok := val.(*userTypes.User)

	if !ok {
		return nil
	}

	return user
}

func UserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		headerParts := strings.Split(authHeader, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := headerParts[1]

		claims, err := ParseToken(token)

		if err != nil {
			http.Error(res, "Unauthorized", http.StatusUnauthorized)
			return
		}

		var user userTypes.User

		user.Email = claims["email"].(string)
		user.ID = uint(claims["id"].(float64))

		ctx := req.Context()
		ctx = WithUser(ctx, &user)
		req = req.WithContext(ctx)

		next.ServeHTTP(res, req)
	})
}

func ParseToken(tok string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Unauthorized")
	}

	if !token.Valid {
		return nil, fmt.Errorf("Unauthorized")
	}

	_, err = token.Claims.GetExpirationTime()
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("Unauthorized")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok {
		return nil, fmt.Errorf("Unauthorized")
	}

	return claims, nil
}
