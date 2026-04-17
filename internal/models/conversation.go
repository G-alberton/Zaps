package models

import "time"

type Conversation struct {
	ID            string    `json:"id"`
	Contact       string    `json: "contact"`
	Status        string    `json: "status"`
	LastMessageAt time.Time `json: "last_message_at"`
	CreatedAt     time.Time `json: "created_at"`
}
