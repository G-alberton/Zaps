package handlers

import (
	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"ZAPS/internal/websocket"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
)

type SendMessageRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
}

func SendMessage(
	mediaService *services.MediaService,
	messageService *services.MessageService,
	conversationService *services.ConversationService,
	hub *websocket.Hub,
) http.HandlerFunc {
	var phoneRegex = regexp.MustCompile(`^\d{10,15}$`)
	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()
		ctx := r.Context()

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SendMessageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if !phoneRegex.MatchString(req.To) {
			http.Error(w, "numero invalido (ex: 5511999999999)", http.StatusBadRequest)
			return
		}

		if req.Body == "" {
			http.Error(w, "mensagem vazia", http.StatusBadRequest)
			return
		}

		log.Println("[SEND] Para:", req.To, "| Msg:", req.Body)

		conversationID, err := conversationService.GetOrCreate(req.To)
		if err != nil {
			http.Error(w, "erro ao obter conversa", 500)
			return
		}

		message := models.Message{
			ID:             uuid.New().String(),
			From:           "system",
			ConversationID: conversationID,
			Type:           "text",
			Body:           req.Body,
			Direction:      "outbound",
			Status:         "pending",
			Timestamp:      time.Now(),
		}

		if err := messageService.SaveMessage(message); err != nil {
			http.Error(w, "erro ao salvar mensagem", 500)
			return
		}

		if hub != nil {
			msgJSON, _ := json.Marshal(message)

			select {
			case hub.Broadcast <- websocket.MessagePayload{
				ConversationID: message.ConversationID,
				Data:           msgJSON,
			}:
			default:
				log.Println("Broadcast cheio")
			}
		}

		go func(msg models.Message) {

			err := mediaService.SendTextMessage(ctx, req.To, req.Body)

			var newStatus string
			if err != nil {
				log.Println("erro envio:", err)
				newStatus = "failed"
			} else {
				newStatus = "sent"
			}

			if err := messageService.UpdateStatus(msg.ID, newStatus); err != nil {
				log.Println("erro ao atualizar status:", err)
			}

			msg.Status = newStatus

			if hub != nil {
				msgJSON, _ := json.Marshal(msg)

				select {
				case hub.Broadcast <- websocket.MessagePayload{
					ConversationID: msg.ConversationID,
					Data:           msgJSON,
				}:
				default:
					log.Println("Broadcast cheio")
				}
			}

		}(message)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":          "processing",
			"conversation_id": conversationID,
		})
	}
}
