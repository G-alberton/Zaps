package webhook

import "log"

func HandleText(msg Message) {
	if msg.Text != nil {
		log.Println("Texto recebido", msg.Text.Body)
	}
}

func HandleImage(msg Message) {
	if msg.Image != nil {
		log.Println("Imagem recebida", msg.Image.ID)
	}
}

func HandleAudio(msg Message) {
	if msg.Audio != nil {
		log.Println("Audio recebido", msg.Audio.ID)
	}
}

func HandleDocument(msg Message) {
	if msg.Document != nil {
		log.Println("Documento recebido", msg.Document.ID)
	}
}
