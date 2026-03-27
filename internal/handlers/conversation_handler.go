package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
)

func GetConversations(conversationService *services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversations := conversationService.List()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversations)
	}
}
