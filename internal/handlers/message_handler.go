package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
)

type MessageResponse struct {
	ID        string `json:"id"`
	From      string `json:"from"`
	Body      string `json:"body"`
	Type      string `json:"type"`
	Direction string `json:"direction"`
	Timestamp string `json:"timestamp"`
}

func GetMessages(messageService *services.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversationID := r.URL.Query().Get("conversation_id")
		if conversationID == "" {
			http.Error(w, "conversation_id is requerid", http.StatusBadRequest)
			return
		}

		messages := messageService.GetByConversation(conversationID)

		var response []MessageResponse

		for _, msg := range messages {
			response = append(response, MessageResponse{
				ID:        msg.ID,
				From:      msg.From,
				Body:      msg.Body,
				Type:      msg.Type,
				Direction: msg.Direction,
				Timestamp: msg.Timestamp.Format("2006-01-02 15:04:05"),
			})
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}

func MarkAsRead(messageService *services.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversationID := r.URL.Query().Get("conversation_id")
		if conversationID == "" {
			http.Error(w, "conversation_id is required", http.StatusBadRequest)
			return
		}

		messageService.MarkAsRead(conversationID)

		w.Write([]byte("ok"))
	}
}

/*func GetMessages(messageService *services.MessageService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversationID := r.URL.Query().Get("conversation_id")
		if conversationID == "" {
			http.Error(w, "conversation_id required", http.StatusBadRequest)
			return
		}

		messages := messageService.GetByConversation(conversationID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(messages)
	}
}*/
