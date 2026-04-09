package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func ServerWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conversationID := r.URL.Query().Get("conversation_id")
	if conversationID == "" {
		http.Error(w, "conversation_id required", http.StatusBadRequest)
		return
	}

	log.Printf("🔌 Nova conexão WS - conversa: %s", conversationID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Erro ao fazer upgrade:", err)
		return
	}

	client := &Client{
		conn:           conn,
		send:           make(chan []byte, 256),
		hub:            hub,
		conversationID: conversationID,
	}

	hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}
