package websocket

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func websocketHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	client := &Client{
		Conn: conn,
		Send: make(chan []byte),
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
		for msg := range client.Send {
			conn.WriteMessage(websocket.TextMessage, msg)
		}
	}()
}
