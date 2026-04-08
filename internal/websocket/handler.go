package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func websocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{
		conn: conn,
		send: make(chan []byte),
	}

	hub.Register <- client

	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				hub.Unregister <- client
				return
			}
			hub.broadcast <- msg
		}
	}()

	go func() {
		for msg := range client.send {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}
