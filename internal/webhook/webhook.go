package webhook

import (
	"ZAPS/internal/models"
	"ZAPS/internal/queue"
	"ZAPS/internal/services"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

const verifyToken = "123456"

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
	q queue.JobQueue,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {

		case http.MethodGet:
			verifyWebhook(w, r)
			return

		case http.MethodPost:
			defer r.Body.Close()

			ctx := r.Context()

			r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

			var event Event

			if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
				log.Println("Erro ao decodificar JSON:", err)
				http.Error(w, "invalid json", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)

			q.Add(queue.High, func() error {
				defer func() {
					if r := recover(); r != nil {
						log.Println("panic no webhook", r)
					}
				}()

				processEvent(
					ctx,
					event,
					contactService,
					messageService,
					mediaService,
					conversationService,
				)
				return nil
			})

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func processEvent(
	ctx context.Context,
	event Event,
	contactService *services.ContactService,
	messageService *services.MessageService,
	mediaService *services.MediaService,
	conversationService *services.ConversationService,
) {

	if len(event.Entry) == 0 {
		log.Println("Webhook vazio (sem entry)")
		return
	}

	for _, entry := range event.Entry {

		for _, change := range entry.Changes {

			if change.Value.Messages == nil {
				continue
			}

			contactsMap := map[string]string{}

			for _, c := range change.Value.Contacts {
				contactsMap[c.WaID] = c.Profile.Name
			}

			for _, msg := range change.Value.Messages {

				processMessage(
					ctx,
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
	ctx context.Context,
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

	log.Printf("Recebido de: %s (%s)", msg.From, name)

	if err := contactService.SaveContact(msg.From, name); err != nil {
		log.Printf("Erro ao salvar contato (%s): %v", msg.From, err)
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
		filePath, err := mediaService.DownloadByID(ctx, mediaID, msg.Type)
		if err != nil {
			log.Printf("Erro ao baixar mídia (%s): %v", mediaID, err)
		} else {
			log.Printf("Mídia salva em: %s", filePath)
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
		log.Printf("Erro ao salvar mensagem (%s): %v", msg.ID, err)
	}
}
