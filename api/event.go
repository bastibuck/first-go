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

	eventId, ok := extractEventID(res, req)
	if !ok {
		return
	}

	event, err := eventHandler.eventStore.GetById(ctx, eventId)
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
	var createEvent types.EventPayloadUpsert

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

func (eventHandler *EventHandler) Update(res http.ResponseWriter, req *http.Request) {
	var updateEvent types.EventPayloadUpsert

	ctx := req.Context()

	eventId, ok := extractEventID(res, req)
	if !ok {
		return
	}

	err := json.NewDecoder(req.Body).Decode(&updateEvent)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Update", http.StatusBadRequest)
		return
	}

	err = eventHandler.eventStore.UpdateEvent(ctx, eventId, &updateEvent)
	if err != nil {
		fmt.Println(err)

		if err == sql.ErrNoRows {
			http.Error(res, "Event not found", http.StatusNotFound)
			return
		}

		http.Error(res, "Something went wrong in Events/Update", http.StatusBadRequest)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (eventHandler *EventHandler) DeleteById(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	eventId, ok := extractEventID(res, req)
	if !ok {
		return
	}

	err := eventHandler.eventStore.DeleteById(ctx, eventId)
	if err != nil {
		fmt.Println(err)

		if err == sql.ErrNoRows {
			http.Error(res, "Event not found", http.StatusNotFound)
			return
		}

		http.Error(res, "Something went wrong in Events/DeleteById", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

// Utils

func extractEventID(res http.ResponseWriter, req *http.Request) (int, bool) {
	idStr := chi.URLParam(req, "id")
	id, err := strconv.Atoi(idStr)

	fmt.Println(id)

	ok := true

	if err != nil {
		http.Error(res, "Invalid event ID", http.StatusBadRequest)
		ok = false
	}

	return id, ok
}
