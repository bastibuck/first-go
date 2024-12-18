package routes

import (
	"first-go/api"

	"github.com/go-chi/chi/v5"
)

func SetupUserRoutes(userHandler *api.UserHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Post("/register", userHandler.RegisterUser)
	router.Post("/login", userHandler.LoginUser)

	return router
}
