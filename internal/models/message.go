package models

//o message serve para salvar a mensagem
type Message struct {
	From    string
	Type    string
	Body    string
	MediaID string
}
