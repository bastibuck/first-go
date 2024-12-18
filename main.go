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

	// Stores
	eventStore := db.NewEventStore(database)
	userStore := db.NewUserStore(database)

	// Handlers
	eventHandler := api.NewEventHandler(eventStore)
	userHandler := api.NewUserHandler(userStore)

	// Routing
	router := chi.NewRouter()

	router.Get("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Welcome to the Event Management API"))
	})

	router.Mount("/events", routes.SetupEventRoutes(eventHandler, eventStore))
	router.Mount("/user", routes.SetupUserRoutes(userHandler))

	// Start the server
	http.ListenAndServe(":8081", router)
}
