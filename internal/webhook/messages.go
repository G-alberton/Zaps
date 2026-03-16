package webhook

import "log"

func handleMessages(event Event) {
	for _, entry := range event.Entry {
		for _, change := range entry.Changes {
			for _, msg := range change.Value.Messages {

				log.Printf("Cliente: %s", msg.From)

				switch msg.Type {

				case "text":
					if msg.Text != nil {
						log.Println("Texto:", msg.Text.Body)
					}
					log.Println("TEXTO:", msg.Text.Body)

				case "imagem":
					log.Println("Imagem ID:", msg.Image.ID)

				case "audio":
					log.Println("Audio ID:", msg.Audio.ID)

				case "document":
					log.Println("Documento ID", msg.Document.ID)

				default:
					log.Println("Tipo desconhecido:", msg.Type)
				}

			}
		}
	}
}
