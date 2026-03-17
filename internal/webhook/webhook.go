package webhook

import (
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

func HandleWebhook(contactService *services.contactService) http.HandlerFunc {

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
					for _, msg := range change.Value.Messages {

						log.Println("Recebido de:", msg.From)

						contactService.SaveContact(msg.From)
					}
				}
			}

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}
