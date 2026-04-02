package handlers

import (
	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func SendMedia(
	mediaService *services.MediaService,
	messageService *services.MessageService,
	conversationService *services.ConversationService,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "error ao ler arquivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		to := r.FormValue("to")
		caption := r.FormValue("caption")

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
		filepath := filepath.Join("uploads/images", filename)

		out, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "erro ao salvar arquivo", 500)
			return
		}
		defer out.Close()

		io.Copy(out, file)

		publicURL := fmt.Sprintf("http://localhost:8080/uploads/images/%s", filename)

		err = mediaService.SendImageByURL(to, publicURL, caption)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		conversationID := conversationService.GetOrCreate(to)

		msg := models.Message{
			From:           to,
			ConversationID: conversationID,
			Type:           "image",
			Body:           caption,
			MediaID:        publicURL,
			Direction:      "outbound",
			Timestamp:      time.Now(),
		}

		messageService.SaveMessage(msg)

		w.Write([]byte("ok"))

	}
}
