package webhook

import (
	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"encoding/json"
	"log"
	"net/http"
)

const verifyToken = "123456"

func verifyWebhook(w http.ResponseWriter, r *http.Request) {

	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "subscribe" && token == verifyToken {
		log.Println("Webhook verificado")
		w.Write([]byte(challenge))
		return
	}

	http.Error(w, "forbidden", http.StatusForbidden)
}

func HandleWebhook(contactService *services.ContactService,
	messageService *services.MessageService,
	mediaService *services.MediaService,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			verifyWebhook(w, r)
			return

		case http.MethodPost:
			defer r.Body.Close()

			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

			var event Event

			err := json.NewDecoder(r.Body).Decode(&event)
			if err != nil {
				log.Println("Erro ao decodificar JSON:", err)
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)

			for _, entry := range event.Entry {
				for _, change := range entry.Changes {

					contractsMap := make(map[string]string)

					for _, c := range change.Value.Contacts {
						contractsMap[c.WaID] = c.Profile.Name
					}
					for _, msg := range change.Value.Messages {

						name := contractsMap[msg.From]

						log.Println("Recebido de:", msg.From, name)

						err := contactService.SaveContact(msg.From, name)
						if err != nil {
							log.Println("Erro ao salvar contato:", err)
						}

						var body string
						var mediaID string

						switch msg.Type {
						case "text":
							if msg.Text != nil {
								body = msg.Text.Body
							}
						case "image":
							if msg.Image != nil {
								mediaID = msg.Image.ID
							}
						case "audio":
							if msg.Audio != nil {
								mediaID = msg.Audio.ID
							}
						case "document":
							if msg.Document != nil {
								mediaID = msg.Document.ID
							}
						}

						if mediaID != "" {
							filePath := "downloads/" + mediaID

							err := mediaService.DowloadByID(mediaID, filePath)
							if err != nil {
								log.Println("Erro ao baixar mídia:", err)
							} else {
								log.Println("Mídia salva em:", filePath)
							}
						}

						message := models.Message{
							From:    msg.From,
							Type:    msg.Type,
							Body:    body,
							MediaID: mediaID,
						}

						err = messageService.SaveMessage(message)
						if err != nil {
							log.Println("Erro ao salvar mensagem", err)
						}
					}
				}
			}

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
