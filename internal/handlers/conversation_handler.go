package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
	"sort"
)

type conversationResponse struct {
	ConversationID string `json:"conversation_id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	LastMessage    string `json:"last_message"`
	Direction      string `json:"direction"`
	Timestamp      string `json:"timestamp"`
	UnreadCount    int    `json:"unread_count"`
}

func GetConversations(
	conversationService *services.ConversationService,
	messageService *services.MessageService,
	contactService *services.ContactService,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		conversations := conversationService.List()

		type temp struct {
			conversationResponse
			rawTime int64
		}

		var tempList []temp

		for _, conv := range conversations {
			lastMsg := messageService.GetLastMessage(conv.ID)

			name := contactService.GetName(conv.Contact)
			if name == "" {
				name = conv.Contact
			}

			resp := conversationResponse{
				ConversationID: conv.ID,
				Name:           name,
				Phone:          conv.Contact,
				UnreadCount:    messageService.CountUnread(conv.ID),
			}

			var rawTime int64

			if lastMsg != nil {
				switch lastMsg.Type {
				case "text":
					resp.LastMessage = lastMsg.Body
				case "image":
					resp.LastMessage = "📷 Imagem"
				case "audio":
					resp.LastMessage = "🎧 Áudio"
				case "document":
					resp.LastMessage = "📄 Documento"
				default:
					resp.LastMessage = "Mensagem"
				}

				resp.Direction = lastMsg.Direction
				resp.Timestamp = lastMsg.Timestamp.Format("15:04")
				rawTime = lastMsg.Timestamp.Unix()
			}

			tempList = append(tempList, temp{
				conversationResponse: resp,
				rawTime:              rawTime,
			})
		}

		sort.Slice(tempList, func(i, j int) bool {
			return tempList[i].rawTime > tempList[j].rawTime
		})

		var response []conversationResponse
		for _, t := range tempList {
			response = append(response, t.conversationResponse)
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		json.NewEncoder(w).Encode(response)
	}
}
