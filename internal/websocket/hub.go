package websocket

import "sync"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	mu         sync.Mutex
	Register   chan *Client
	Unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		clients:   make(map[*Client]bool),
		broadcast: make(chan []byte),
	}
}

func (h *Hub) Run() {
	for {
		msg := <-h.broadcast

		h.mu.Lock()
		for client := range h.clients {
			select {
			case client.send <- msg:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
		h.mu.Unlock()
	}
}
