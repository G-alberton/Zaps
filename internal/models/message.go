package models

import "time"

//o message serve para salvar a mensagem
type Message struct {
	From      string
	Type      string
	Body      string
	MediaID   string
	Timestamp time.Time
}
