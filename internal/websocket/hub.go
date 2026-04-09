package websocket

import "sync"

type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan MessagePayload
	mu         sync.Mutex
	Register   chan *Client
	Unregister chan *Client
}

type MessagePayload struct {
	ConversationID string
	Data           []byte
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan MessagePayload, 256),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.Register:
			h.mu.Lock()

			if h.Rooms[client.conversationID] == nil {
				h.Rooms[client.conversationID] = make(map[*Client]bool)
			}

			h.Rooms[client.conversationID][client] = true

			h.mu.Unlock()

		case client := <-h.Unregister:
			h.mu.Lock()

			if clients, ok := h.Rooms[client.conversationID]; ok {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.send)
				}
			}

			h.mu.Unlock()

		case msg := <-h.Broadcast:
			h.mu.Lock()

			if clients, ok := h.Rooms[msg.ConversationID]; ok {
				for client := range clients {
					select {
					case client.send <- msg.Data:
					default:
						close(client.send)
						delete(clients, client)
					}
				}
			}

			h.mu.Unlock()
		}
	}
}
