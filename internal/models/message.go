package models

import "time"

type Message struct {
	ID             string
	ConversationID string
	From           string
	Type           string
	Body           string
	MediaID        string
	MediaURL       string
	Timestamp      time.Time
	Direction      string
	Status         string
	Read           bool
}
