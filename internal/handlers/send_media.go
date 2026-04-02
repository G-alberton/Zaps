package handlers

import (
	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

		contentType := header.Header.Get("Content-type")

		to := r.FormValue("to")
		caption := r.FormValue("caption")

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)

		var folder string

		if strings.HasPrefix(contentType, "image/") {
			folder = "uploads/images"
		} else if strings.HasPrefix(contentType, "audio/") {
			folder = "uploads/audio"
		} else {
			folder = "uploads/files"
		}
		filepath := filepath.Join(folder, filename)

		out, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "erro ao salvar arquivo", 500)
			return
		}
		defer out.Close()

		io.Copy(out, file)

		publicURL := fmt.Sprintf("http://localhost:8080/uploads/images/%s", filename)

		if strings.HasPrefix(contentType, "image/") {
			err = mediaService.SendImageByURL(to, publicURL, caption)
		} else if strings.HasPrefix(contentType, "audio/") {
			mediaID, errUpload := mediaService.UploadMedia(filepath)
			if errUpload != nil {
				http.Error(w, errUpload.Error(), 500)
				return
			}

			err = mediaService.SendAudioByID(to, mediaID)
		}

		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		msgType := "file"

		if strings.HasPrefix(contentType, "image/") {
			msgType = "image"
		} else if strings.HasPrefix(contentType, "audio/") {
			msgType = "audio"
		}

		conversationID := conversationService.GetOrCreate(to)

		msg := models.Message{
			From:           to,
			ConversationID: conversationID,
			Type:           msgType,
			Body:           caption,
			MediaID:        publicURL,
			Direction:      "outbound",
			Timestamp:      time.Now(),
		}

		os.MkdirAll(folder, os.ModePerm)

		messageService.SaveMessage(msg)

		w.Write([]byte("ok"))

	}
}
