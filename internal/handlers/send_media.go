package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
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
	messageService *services.MessageService,
	conversationService *services.ConversationService,
	sendService *services.SendService,
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
			w.Header().Set("Content-Type", "application/json")
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

		to := strings.TrimSpace(r.FormValue("to"))
		caption := strings.TrimSpace(r.FormValue("caption"))

		if !phoneRegex.MatchString(to) {
			http.Error(w, "numero invalido", http.StatusBadRequest)
			return
		}

		buffer := make([]byte, 512)
		n, _ := file.Read(buffer)
		realType := http.DetectContentType(buffer[:n])
		file.Seek(0, io.SeekStart)

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

		safeName := fileSafeRegex.ReplaceAllString(header.Filename, "")
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), safeName)
		fullPath := filepath.Join("uploads", folder, filename)

		os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)

		out, _ := os.Create(fullPath)
		defer out.Close()
		io.Copy(out, file)

		baseURL := os.Getenv("BASE_URL")
		publicURL := fmt.Sprintf("%s/uploads/%s/%s", baseURL, folder, filename)

		conversationID, _ := conversationService.GetOrCreate(to)

		message := models.Message{
			ID:             messageID,
			ConversationID: conversationID,
			Type:           msgType,
			Body:           caption,
			MediaURL:       publicURL,
			Direction:      "outbound",
			Status:         "pending",
			Timestamp:      time.Now(),
		}

		messageService.SaveMessage(message)

		if hub != nil {
			broadcast(hub, message)
		}

		go sendService.SendWithRetry(
			message,
			to,
			caption,
			fullPath,
			filename,
		)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":     "processing",
			"message_id": messageID,
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
	}
}
