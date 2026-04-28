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

	return func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		ctx := r.Context()

		var phoneRegex = regexp.MustCompile(`^\d{10,15}$`)

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
			http.Error(w, "numero invalido (use formato internacional: 5511999999999)", http.StatusBadRequest)
			return
		}

		log.Println("[SEND] Para:", req.To, "| Msg:", req.Body)

		conversationID, err := conversationService.GetOrCreate(req.To)
		if err != nil {
			http.Error(w, "erro ao obter conversa", 500)
			return
		}

		message := models.Message{
			From:           "system",
			ConversationID: conversationID,
			Type:           "text",
			Body:           req.Body,
			Direction:      "outbound",
			Status:         "pending",
			Timestamp:      time.Now(),
		}

		if err := messageService.SaveMessage(message); err != nil {
			log.Println("erro ao salvar mensagem:", err)
		}

		err = mediaService.SendTextMessage(ctx, req.To, req.Body)

		if err != nil {
			message.Status = "failed"
			log.Println("erro ao enviar mensagem:", err)
		} else {
			message.Status = "sent"
		}

		if err := messageService.UpdateStatus(message.ID, message.Status); err != nil {
			log.Println("erro ao atualizar status:", err)
		}

		if hub != nil {
			msgJSON, err := json.Marshal(message)
			if err == nil {
				select {
				case hub.Broadcast <- websocket.MessagePayload{
					ConversationID: message.ConversationID,
					Data:           msgJSON,
				}:
				default:
					log.Println("Broadcast cheio, descartando mensagem")
				}
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":          message.Status,
			"conversation_id": conversationID,
			"to":              req.To,
			"body":            req.Body,
		})
	}

}
