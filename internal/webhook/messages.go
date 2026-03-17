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
					HandleText(msg)

				case "imagem":
					HandleImage(msg)

				case "audio":
					HandleAudio(msg)

				case "document":
					HandleDocument(msg)

				default:
					log.Println("Tipo desconhecido:", msg.Type)
				}

			}
		}
	}
}
