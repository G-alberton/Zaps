package handlers

import (
	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"ZAPS/internal/websocket"
	"encoding/json"
	"log"
	"net/http"
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

		ctx := r.Context()

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req SendMessageRequest

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}

		if req.To == "" || req.Body == "" {
			http.Error(w, "to and body requerid", http.StatusBadRequest)
			return
		}

		log.Println("[SEND] Para:", req.To, "| Msg:", req.Body)

		err = mediaService.SendTextMessage(ctx, req.To, req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		conversationID, err := conversationService.GetOrCreate(req.To)
		if err != nil {
			http.Error(w, "erro ao obter conversa", 500)
			return
		}

		message := models.Message{
			From:           req.To,
			ConversationID: conversationID,
			Type:           "text",
			Body:           req.Body,
			Direction:      "outbound",
			Timestamp:      time.Now(),
		}

		err = messageService.SaveMessage(message)
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
		if err != nil {
			log.Println("erro ao salvar mensagem enviada:", err)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":          "ok",
			"conversation_id": conversationID,
			"to":              req.To,
			"body":            req.Body,
		})
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"status":          "ok",
			"conversation_id": conversationID,
			"to":              req.To,
			"body":            req.Body,
		}); err != nil {
			log.Println("erro ao retornar resposta:", err)
		}
	}

}
