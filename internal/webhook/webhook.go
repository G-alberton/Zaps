package webhook

import (
	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const verifyToken = "123456"

//var verifyToken = os.Getenv("VERIFY_TOKEN") //colocar na .env

func init() {
	if verifyToken == "" {
		log.Fatal("VERIFY_TOKEN não definido")
	}
}

func verifyWebhook(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("hub.mode")
	token := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if mode == "" || token == "" {
		http.Error(w, "missing params", http.StatusBadRequest)
		return
	}

	if mode == "subscribe" && token == verifyToken {
		log.Println("Webhook verificado com sucesso")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(challenge))
		return
	}

	log.Println("Falha na verificação do webhook")
	http.Error(w, "forbidden", http.StatusForbidden)
}

func HandleWebhook(
	contactService *services.ContactService,
	messageService *services.MessageService,
	mediaService *services.MediaService,
	conversationService *services.ConversationService,
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

			if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
				log.Println("Erro ao decodificar JSON:", err)
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)

			go processEvent(event, contactService, messageService, mediaService, conversationService)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func processEvent(
	event Event,
	contactService *services.ContactService,
	messageService *services.MessageService,
	mediaService *services.MediaService,
	conversationService *services.ConversationService,
) {

	for _, entry := range event.Entry {
		for _, change := range entry.Changes {

			contactsMap := map[string]string{}

			for _, c := range change.Value.Contacts {
				contactsMap[c.WaID] = c.Profile.Name
			}

			for _, msg := range change.Value.Messages {
				processMessage(
					msg,
					contactsMap,
					contactService,
					messageService,
					mediaService,
					conversationService,
				)
			}
		}
	}
}

func processMessage(
	msg Message,
	contactsMap map[string]string,
	contactService *services.ContactService,
	messageService *services.MessageService,
	mediaService *services.MediaService,
	conversationService *services.ConversationService,
) {

	name, ok := contactsMap[msg.From]
	if !ok {
		name = "Unknown"
	}
	conversationID := conversationService.GetOrCreate(msg.From)

	log.Println("Recebido de:", msg.From, name)

	if err := contactService.SaveContact(msg.From, name); err != nil {
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
		if filePath, err := mediaService.DownloadByID(mediaID, msg.Type); err != nil {
			log.Println("Erro ao baixar mídia:", err)
		} else {
			log.Println("Mídia salva em:", filePath)
		}
	}

	tsInt, err := strconv.ParseInt(msg.Timestamp, 10, 64)
	if err != nil {
		log.Println("Erro ao converter timestamp:", err)
		tsInt = time.Now().Unix()
	}

	timestamp := time.Unix(tsInt, 0)

	message := models.Message{
		ID:             msg.ID,
		ConversationID: conversationID,
		From:           msg.From,
		Type:           msg.Type,
		Body:           body,
		MediaID:        mediaID,
		Timestamp:      timestamp,
		Direction:      "inbound",
	}

	if err := messageService.SaveMessage(message); err != nil {
		log.Println("Erro ao salvar mensagem:", err)
	}
}
