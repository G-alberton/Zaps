package models

import "time"

type Message struct {
	ID             string
	ConversationID string
	From           string
	Type           string
	Body           string
	MediaID        string
	Timestamp      time.Time
	Direction      string
}
