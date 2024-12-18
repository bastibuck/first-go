package routes

import (
	"first-go/api"
	"first-go/db"
	"first-go/routes/middleware"

	"github.com/go-chi/chi/v5"
)

func SetupEventRoutes(eventHandler *api.EventHandler, eventStore db.EventStore) *chi.Mux {
	router := chi.NewRouter()

	// public routes
	router.Get("/", eventHandler.GetAll)
	router.With(middleware.EventExistence(eventStore)).Get("/{id}", eventHandler.GetById)

	// protected routes
	router.With(middleware.UserAuthentication).Post("/", eventHandler.Create)
	router.With(middleware.UserAuthentication).With(middleware.EventExistence(eventStore)).With(middleware.UserEvent(eventStore)).Put("/{id}", eventHandler.Update)
	router.With(middleware.UserAuthentication).With(middleware.EventExistence(eventStore)).With(middleware.UserEvent(eventStore)).Delete("/{id}", eventHandler.DeleteById)

	return router
}
