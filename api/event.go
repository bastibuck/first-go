package api

import (
	"database/sql"
	"encoding/json"
	"first-go/db"
	"first-go/types"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type EventHandler struct {
	eventStore db.EventStore
}

func NewEventHandler(eventStore db.EventStore) *EventHandler {
	return &EventHandler{
		eventStore,
	}
}

func (eventHandler *EventHandler) GetAll(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	events, err := eventHandler.eventStore.GetAll(ctx)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/GetAll", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(events)
}

func (eventHandler *EventHandler) GetById(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	idStr := chi.URLParam(req, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(res, "Invalid event ID", http.StatusBadRequest)
		return
	}

	event, err := eventHandler.eventStore.GetById(ctx, id)
	if err != nil {
		fmt.Println(err)

		if err == sql.ErrNoRows {
			http.Error(res, "Event not found", http.StatusNotFound)
			return
		}

		http.Error(res, "Something went wrong in Events/GetById", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(event)
}

func (eventHandler *EventHandler) Create(res http.ResponseWriter, req *http.Request) {
	var createEvent types.CreateEvent

	ctx := req.Context()

	err := json.NewDecoder(req.Body).Decode(&createEvent)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Create", http.StatusBadRequest)
		return
	}

	err = eventHandler.eventStore.AddEvent(ctx, &createEvent)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Create", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusCreated)
}
