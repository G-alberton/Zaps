package handlers

import (
	"ZAPS/internal/services"
	"encoding/json"
	"net/http"
)

type SendMessageRequest struct {
	To   string `json:"to"`
	Body string `json:"body"`
}

func SendMessage(mediaService *services.MediaService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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

		err = mediaService.SendTextMessage(req.To, req.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write([]byte("mensagem enviada"))
	}
}
