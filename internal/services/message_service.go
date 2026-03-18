package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
	"log"
	"time"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SaveMessage(msg models.Message) error {
	if msg.From == "" {
		return nil
	}

	switch msg.Type {
	case "text":
		if msg.Body == "" {
			log.Println("Mensagem de Texto Vazio")
			return nil
		}
	case "image", "audio", "document":
		if msg.MediaID == "" {
			log.Println("Mídia sem ID")
			return nil
		}
	default:
		log.Println("Tipo não suportado:", msg.Type)
		return nil
	}

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	err := s.repo.Save(msg)
	if err != nil {
		log.Println("Erro ao salvar mensagem:", err)
		return err
	}

	log.Println("Mensagem salva:", msg.From, msg.Type)
	return nil
}
