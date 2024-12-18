package main

import (
	"first-go/api"
	"first-go/db"
	"first-go/routes"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	fmt.Println("Event management is starting...")

	database := db.OpenConnection()

	eventStore := db.NewEventStore(database)
	eventHandler := api.NewEventHandler(eventStore)

	// Create main router
	router := chi.NewRouter()

	// Mount sub-routers
	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Welcome to the Event Management API"))
	})

	router.Mount("/events", routes.SetupEventRoutes(eventHandler))

	http.ListenAndServe(":8081", router)
}
