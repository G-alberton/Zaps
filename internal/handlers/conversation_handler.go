package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
	"sort"
)

/*func GetConversations(conversationService *services.ConversationService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversations := conversationService.List()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(conversations)
	}
}*/

type conversationResponse struct {
	ConversationID string `json:"conversation_id"`
	Phone          string `json:"phone"`
	LastMessage    string `json:"last_message"`
	Direction      string `json:"direction"`
	Timestamp      string `json:"timestamp"`
}

func GetConversations(
	conversationService *services.ConversationService,
	messageService *services.MessageService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversations := conversationService.List()

		var response []conversationResponse

		for _, conv := range conversations {
			lastMsg := messageService.GetLastMessage(conv.ID)

			resp := conversationResponse{
				ConversationID: conv.ID,
				Phone:          conv.Contact,
			}

			if lastMsg != nil {
				if lastMsg.Type == "text" {
					resp.LastMessage = lastMsg.Body
				} else {
					resp.LastMessage = "📎 " + lastMsg.Type
				}
				resp.Direction = lastMsg.Direction
				resp.Timestamp = lastMsg.Timestamp.Format("2006-01-02 15:04:05")
			}

			response = append(response, resp)
		}

		sort.Slice(response, func(i, j int) bool {
			return response[i].Timestamp > response[j].Timestamp
		})

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}
