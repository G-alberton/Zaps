package services

import (
	"ZAPS/internal/models"
	"sync"

	"github.com/google/uuid"
)

type ConversationService struct {
	conversations map[string]models.Conversation
	mu            sync.Mutex
}

func NewConversationService() *ConversationService {
	return &ConversationService{
		conversations: make(map[string]models.Conversation),
	}
}

func (s *ConversationService) GetOrCreate(contact string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	// se já existe, retorna
	if conv, ok := s.conversations[contact]; ok {
		return conv.ID
	}

	// cria nova
	id := uuid.New().String()

	s.conversations[contact] = models.Conversation{
		ID:      id,
		Contact: contact,
	}

	return id
}
