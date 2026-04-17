package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
	"sort"
)

type ConversationHandler struct {
	service        *services.ConversationService
	messageService *services.MessageService
	contactService *services.ContactService
}

type conversationResponse struct {
	ConversationID string `json:"conversation_id"`
	Name           string `json:"name"`
	Phone          string `json:"phone"`
	LastMessage    string `json:"last_message"`
	Direction      string `json:"direction"`
	Timestamp      string `json:"timestamp"`
	UnreadCount    int    `json:"unread_count"`
}

func (h *ConversationHandler) GetConversations(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	conversations, err := h.service.List()
	if err != nil {
		http.Error(w, "erro ao buscar conversas", 500)
		return
	}

	type temp struct {
		conversationResponse
		rawTime int64
	}

	var tempList []temp

	for _, conv := range conversations {

		lastMsg := h.messageService.GetLastMessage(conv.ID)

		name := h.contactService.GetName(conv.Contact)
		if name == "" {
			name = conv.Contact
		}

		resp := conversationResponse{
			ConversationID: conv.ID,
			Name:           name,
			Phone:          conv.Contact,
			UnreadCount:    h.messageService.CountUnread(conv.ID),
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *ConversationHandler) GetOrCreate(w http.ResponseWriter, r *http.Request) {
	contact := r.URL.Query().Get("contact")
	if contact == "" {
		http.Error(w, "contact obrigatorio", 400)
		return
	}

	id, err := h.service.GetOrCreate(contact)
	if err != nil {
		http.Error(w.err.Error(), 500)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"conversation_id": id,
	})
}
