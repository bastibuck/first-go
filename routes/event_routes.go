package routes

import (
	"first-go/api"

	"github.com/go-chi/chi/v5"
)

func SetupEventRoutes(eventHandler *api.EventHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Get("/", eventHandler.GetAll)
	router.Get("/{id}", eventHandler.GetById)
	router.Post("/", eventHandler.Create)
	router.Delete("/{id}", eventHandler.DeleteById)

	return router
}
