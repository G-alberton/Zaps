package webhook

import (
	"ZAPS/internal/models"
	"ZAPS/internal/queue"
	"ZAPS/internal/services"
	"ZAPS/internal/websocket"
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
	hub *websocket.Hub,
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

			// ✅ resposta imediata (recomendado pelo Meta)
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("EVENT_RECEIVED"))

			q.Add(queue.High, func() error {

				jobCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()

				processEvent(
					jobCtx,
					event,
					contactService,
					messageService,
					mediaService,
					conversationService,
					hub,
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
	hub *websocket.Hub,
) {

	if len(event.Entry) == 0 {
		log.Println("Webhook vazio (sem entry)")
		return
	}

	for _, entry := range event.Entry {

		for _, change := range entry.Changes {

			if change.Value.Statuses != nil {
				for _, status := range change.Value.Statuses {

					log.Printf("📊 Status: %s → %s", status.ID, status.Status)

					if err := messageService.UpdateStatus(status.ID, status.Status); err != nil {
						log.Println("erro ao atualizar status:", err)
					}
				}
			}

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
					hub,
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
	hub *websocket.Hub,
) {

	select {
	case <-ctx.Done():
		log.Println("Context cancelado no processMessage")
		return
	default:
	}

	exists, err := messageService.Exists(msg.ID)
	if err == nil && exists {
		log.Println("Mensagem duplicada ignorada:", msg.ID)
		return
	}

	name := contactsMap[msg.From]
	if name == "" {
		name = "Unknown"
	}

	conversationID, err := conversationService.GetOrCreate(msg.From)
	if err != nil {
		log.Printf("Erro ao obter/criar conversa (%s): %v", msg.From, err)
		return
	}

	log.Printf("📩 MsgID: %s | From: %s (%s) | Type: %s",
		msg.ID, msg.From, name, msg.Type)

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
		go func() {
			filePath, err := mediaService.DownloadByID(context.Background(), mediaID, msg.Type)
			if err != nil {
				log.Printf("Erro ao baixar mídia (%s): %v", mediaID, err)
				return
			}
			log.Printf("Mídia salva em: %s", filePath)
		}()
	}

	tsInt, err := strconv.ParseInt(msg.Timestamp, 10, 64)
	if err != nil {
		log.Println("Erro ao converter timestamp:", err)
		tsInt = time.Now().Unix()
	}

	var timestamp time.Time
	if tsInt > 1e12 {
		timestamp = time.UnixMilli(tsInt)
	} else {
		timestamp = time.Unix(tsInt, 0)
	}

	message := models.Message{
		ID:             msg.ID,
		ConversationID: conversationID,
		From:           msg.From,
		Type:           msg.Type,
		Body:           body,
		MediaID:        mediaID,
		Timestamp:      timestamp,
		Direction:      "inbound",
		Status:         "received",
	}

	if err := messageService.SaveMessage(message); err != nil {
		log.Printf("Erro ao salvar mensagem (%s): %v", msg.ID, err)
	}

	if hub != nil {
		msgJSON, err := json.Marshal(message)
		if err != nil {
			log.Println("Erro ao serializar mensagem:", err)
			return
		}

		select {
		case hub.Broadcast <- websocket.MessagePayload{
			ConversationID: conversationID,
			Data:           msgJSON,
		}:
			log.Println("Mensagem enviada para WebSocket")
		default:
			log.Println("Canal Broadcast cheio, descartando mensagem")
		}
	}
}
