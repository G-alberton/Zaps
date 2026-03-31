package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
	"log"
	"sync"
	"time"
)

type MessageService struct {
	repo     *repository.MessageRepository
	messages []models.Message
	mu       sync.Mutex
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{
		repo:     repo,
		messages: []models.Message{},
	}
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

	s.mu.Lock()
	s.messages = append(s.messages, msg)
	s.mu.Unlock()

	if s.repo != nil {
		err := s.repo.Save(msg)
		if err != nil {
			log.Println("Erro ao salvar mensagem:", err)
		}
	}

	log.Println("Mensagem salva:", msg.From, msg.Type, msg.Direction)
	return nil
}

func (s *MessageService) GetByConversation(conversationID string) []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	var result []models.Message

	for _, msg := range s.messages {
		if msg.ConversationID == conversationID {
			result = append(result, msg)
		}
	}

	return result
}

func (s *MessageService) GetLastMessage(conversationID string) *models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	var last *models.Message

	for i := range s.messages {
		msg := &s.messages[i]

		if msg.ConversationID == conversationID {
			if last == nil || msg.Timestamp.After(last.Timestamp) {
				last = msg
			}
		}
	}

	return last
}

func (s *MessageService) GetAll() []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.messages
}
