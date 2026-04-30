package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"ZAPS/internal/models"
	"ZAPS/internal/websocket"
)

type SendService struct {
	MediaService   *MediaService
	MessageService *MessageService
	Hub            *websocket.Hub
}

func NewSendService(
	media *MediaService,
	message *MessageService,
	hub *websocket.Hub,
) *SendService {
	return &SendService{
		MediaService:   media,
		MessageService: message,
		Hub:            hub,
	}
}

func (s *SendService) SendWithRetry(
	msg models.Message,
	to string,
	caption string,
	filePath string,
	filename string,
) {
	maxRetries := 3
	delay := 2 * time.Second

	ctx := context.Background()

	for attempt := 1; attempt <= maxRetries; attempt++ {

		log.Printf("[TENTATIVA %d] msg=%s", attempt, msg.ID)

		mediaID, err := s.MediaService.UploadMedia(ctx, filePath)
		if err != nil {
			log.Println("Erro upload:", err)
			s.waitRetry(attempt, delay)
			continue
		}

		msg.MediaID = mediaID

		err = s.sendByType(ctx, msg, to, caption, filename)
		if err == nil {
			s.updateStatus(msg, "sent")
			return
		}

		log.Println("Erro envio:", err)

		if attempt < maxRetries {
			s.waitRetry(attempt, delay)
		}
	}

	s.updateStatus(msg, "failed")
	log.Println("Falha definitiva:", msg.ID)
}

func (s *SendService) sendByType(
	ctx context.Context,
	msg models.Message,
	to string,
	caption string,
	filename string,
) error {

	switch msg.Type {
	case "image":
		return s.MediaService.SendImageByID(ctx, to, msg.MediaID, caption)
	case "audio":
		return s.MediaService.SendAudioByID(ctx, to, msg.MediaID)
	case "file":
		return s.MediaService.SendDocumentByID(ctx, to, msg.MediaID, caption, filename)
	default:
		return nil
	}
}

func (s *SendService) updateStatus(msg models.Message, status string) {
	err := s.MessageService.UpdateStatus(msg.ID, status)
	if err != nil {
		log.Println("Erro ao atualizar status:", err)
	}

	msg.Status = status

	if s.Hub != nil {
		s.broadcast(msg)
	}
}

func (s *SendService) broadcast(msg models.Message) {
	msgJSON, _ := json.Marshal(msg)

	select {
	case s.Hub.Broadcast <- websocket.MessagePayload{
		ConversationID: msg.ConversationID,
		Data:           msgJSON,
	}:
	default:
		log.Println("Broadcast cheio")
	}
}

func (s *SendService) waitRetry(attempt int, base time.Duration) {
	backoff := time.Duration(1<<attempt) * base // exponencial
	log.Printf("Aguardando %v para retry...", backoff)
	time.Sleep(backoff)
}
