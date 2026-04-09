package websocket

import "sync"

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan []byte
	mu         sync.Mutex
	Register   chan *Client
	Unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[*Client]bool),
		Broadcast:  make(chan []byte, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client] = true
			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.Lock()
			for client := range h.Clients {
				select {
				case client.send <- msg:
				default:
					close(client.send)
					delete(h.Clients, client)
				}
			}
			h.mu.Unlock()
		}
	}
}
