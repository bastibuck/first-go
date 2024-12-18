package routes

import (
	"first-go/api"
	"first-go/routes/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupEventRoutes(eventHandler *api.EventHandler) *chi.Mux {
	router := chi.NewRouter()

	// public routes
	router.Get("/", eventHandler.GetAll)
	router.Get("/{id}", eventHandler.GetById)

	// protected routes
	router.With(middleware.UserAuthentication).Post("/", eventHandler.Create)
	router.With(middleware.UserAuthentication).With(middleware.EventProtection).Put("/{id}", eventHandler.Update)
	router.With(middleware.UserAuthentication).With(middleware.EventProtection).Delete("/{id}", eventHandler.DeleteById)

	return router
}
