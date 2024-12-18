package api

import (
	"encoding/json"
	"first-go/db"
	"first-go/routes/middleware"
	eventTypes "first-go/types/event"
	"first-go/utils"
	"fmt"
	"net/http"
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

	eventId, ok := utils.ExtractEventID(res, req)
	if !ok {
		return
	}

	event, err := eventHandler.eventStore.GetById(ctx, eventId)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/GetById", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(res).Encode(event)
}

func (eventHandler *EventHandler) Create(res http.ResponseWriter, req *http.Request) {
	var createEvent eventTypes.EventPayloadUpsert

	ctx := req.Context()

	user := middleware.User(ctx)

	err := json.NewDecoder(req.Body).Decode(&createEvent)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Create", http.StatusInternalServerError)
		return
	}

	err = eventHandler.eventStore.AddEvent(ctx, &createEvent, user.ID)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Create", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusCreated)
}

func (eventHandler *EventHandler) Update(res http.ResponseWriter, req *http.Request) {
	var updateEvent eventTypes.EventPayloadUpsert

	ctx := req.Context()

	eventId, ok := utils.ExtractEventID(res, req)
	if !ok {
		return
	}

	err := json.NewDecoder(req.Body).Decode(&updateEvent)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Update", http.StatusInternalServerError)
		return
	}

	err = eventHandler.eventStore.UpdateEvent(ctx, eventId, &updateEvent)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/Update", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}

func (eventHandler *EventHandler) DeleteById(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	eventId, ok := utils.ExtractEventID(res, req)
	if !ok {
		return
	}

	err := eventHandler.eventStore.DeleteById(ctx, eventId)
	if err != nil {
		fmt.Println(err)
		http.Error(res, "Something went wrong in Events/DeleteById", http.StatusInternalServerError)
		return
	}

	res.WriteHeader(http.StatusNoContent)
}
