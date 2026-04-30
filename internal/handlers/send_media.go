package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"ZAPS/internal/models"
	"ZAPS/internal/services"
	"ZAPS/internal/websocket"

	"github.com/google/uuid"
)

var phoneRegex = regexp.MustCompile(`^\d{10,15}$`)
var fileSafeRegex = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func SendMedia(
	mediaService *services.MediaService,
	messageService *services.MessageService,
	conversationService *services.ConversationService,
	hub *websocket.Hub,
) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		messageID := r.Header.Get("X-Message-ID")
		if messageID == "" {
			messageID = uuid.New().String()
		}

		if _, err := messageService.GetByID(messageID); err == nil {
			json.NewEncoder(w).Encode(map[string]string{
				"status": "already_processed",
			})
			return
		} else if err != sql.ErrNoRows {
			http.Error(w, "erro interno", 500)
			return
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "erro ao ler arquivo", http.StatusBadRequest)
			return
		}
		defer file.Close()

		if header.Size == 0 {
			http.Error(w, "arquivo vazio", http.StatusBadRequest)
			return
		}

		to := strings.TrimSpace(r.FormValue("to"))
		caption := strings.TrimSpace(r.FormValue("caption"))

		if !phoneRegex.MatchString(to) {
			http.Error(w, "numero invalido (ex: 5511999999999)", http.StatusBadRequest)
			return
		}

		buffer := make([]byte, 512)
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			http.Error(w, "erro ao ler arquivo", 500)
			return
		}

		if n == 0 {
			http.Error(w, "arquivo inválido", 400)
			return
		}

		realType := http.DetectContentType(buffer[:n])

		if !strings.HasPrefix(realType, "image/") &&
			!strings.HasPrefix(realType, "audio/") &&
			realType != "application/pdf" {
			http.Error(w, "tipo de arquivo não permitido", 400)
			return
		}

		_, err = file.Seek(0, io.SeekStart)
		if err != nil {
			http.Error(w, "erro ao reposicionar arquivo", 500)
			return
		}

		var folder, msgType string

		switch {
		case strings.HasPrefix(realType, "image/"):
			folder = "images"
			msgType = "image"
		case strings.HasPrefix(realType, "audio/"):
			folder = "audio"
			msgType = "audio"
		default:
			folder = "files"
			msgType = "file"
		}

		safeName := strings.ReplaceAll(header.Filename, " ", "_")
		safeName = fileSafeRegex.ReplaceAllString(safeName, "")

		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), safeName)
		fullPath := filepath.Join("uploads", folder, filename)

		if err := os.MkdirAll(filepath.Dir(fullPath), os.ModePerm); err != nil {
			http.Error(w, "erro ao criar diretório", 500)
			return
		}

		out, err := os.Create(fullPath)
		if err != nil {
			http.Error(w, "erro ao salvar arquivo", 500)
			return
		}
		defer out.Close()

		if _, err = io.Copy(out, file); err != nil {
			http.Error(w, "erro ao salvar arquivo", 500)
			return
		}

		baseURL := os.Getenv("BASE_URL")
		if baseURL == "" {
			http.Error(w, "BASE_URL não configurada", 500)
			return
		}

		publicURL := fmt.Sprintf("%s/uploads/%s/%s", baseURL, folder, filename)

		conversationID, err := conversationService.GetOrCreate(to)
		if err != nil {
			http.Error(w, "erro ao obter conversa", 500)
			return
		}

		message := models.Message{
			ID:             messageID,
			ConversationID: conversationID,
			From:           "system",
			Type:           msgType,
			Body:           caption,
			MediaURL:       publicURL,
			Direction:      "outbound",
			Status:         "pending",
			Timestamp:      time.Now(),
		}

		if err := messageService.SaveMessage(message); err != nil {
			if strings.Contains(err.Error(), "duplicate") {
				json.NewEncoder(w).Encode(map[string]string{
					"status": "already_processed",
				})
				return
			}
			http.Error(w, err.Error(), 500)
			return
		}

		if hub != nil {
			broadcast(hub, message)
		}

		go func(msg models.Message, filePath string, filename string) {

			ctx := context.Background()

			mediaID, err := mediaService.UploadMedia(ctx, filePath)

			var newStatus string

			if err != nil {
				log.Println("erro upload:", err)
				newStatus = "failed"
			} else {
				msg.MediaID = mediaID

				switch msg.Type {
				case "image":
					err = mediaService.SendImageByID(ctx, to, mediaID, caption)
				case "audio":
					err = mediaService.SendAudioByID(ctx, to, mediaID)
				case "file":
					err = mediaService.SendDocumentByID(ctx, to, mediaID, caption, filename)
				}

				if err != nil {
					log.Println("erro envio:", err)
					newStatus = "failed"
				} else {
					newStatus = "sent"
				}
			}

			if err := messageService.UpdateStatus(msg.ID, newStatus); err != nil {
				log.Println("erro ao atualizar status:", err)
			}

			msg.Status = newStatus

			if hub != nil {
				broadcast(hub, msg)
			}

		}(message, fullPath, safeName)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":          "processing",
			"conversation_id": conversationID,
			"message_id":      messageID,
		})
	}
}

func broadcast(hub *websocket.Hub, msg models.Message) {
	msgJSON, _ := json.Marshal(msg)

	select {
	case hub.Broadcast <- websocket.MessagePayload{
		ConversationID: msg.ConversationID,
		Data:           msgJSON,
	}:
	default:
		log.Println("Broadcast cheio")
	}
}
