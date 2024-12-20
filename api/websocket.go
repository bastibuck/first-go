package api

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan WebSocketMessage)

// all message types
type MessageType string

const (
	EventSignUp MessageType = "event_signup"
)

// message format
type WebSocketMessage struct {
	Type    MessageType `json:"type"`
	Payload Payload     `json:"payload"`
}

type Payload interface {
	isPayload()
}

// payload types
type EventSignUpPayload struct {
	EventID uint `json:"event_id"`
	NewPax  int  `json:"new_pax"`
}

func (EventSignUpPayload) isPayload() {}

func HandleWebSocketConnections(res http.ResponseWriter, req *http.Request) {
	ws, err := upgrader.Upgrade(res, req, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ws.Close()

	clients[ws] = true

	for {
		var msg WebSocketMessage
		err := ws.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Printf("Unexpected close error: %v\n", err)
			} else {
				fmt.Printf("WebSocket closed: %v\n", err)
			}
			delete(clients, ws)
			break
		}
	}
}

func HandleWebSocketMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				fmt.Println(err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
