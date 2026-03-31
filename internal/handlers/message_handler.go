package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
)

func GetMessages(messageService *services.MessageService) http.HandlerFunc {
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
}
