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

func (s *ConversationService) GetOrCreate(contact string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if conv, ok := s.conversations[contact]; ok {
		return conv.ID, nil
	}

	id := uuid.New().String()

	s.conversations[contact] = models.Conversation{
		ID:      id,
		Contact: contact,
	}

	return id, nil
}

func (s *ConversationService) List() []models.Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	conversations := []models.Conversation{}

	for _, conv := range s.conversations {
		conversations = append(conversations, conv)
	}

	return conversations
}

func (s *ConversationService) GetAll() []models.Conversation {
	s.mu.Lock()
	defer s.mu.Unlock()

	var list []models.Conversation

	for _, c := range s.conversations {
		list = append(list, c)
	}

	return list
}
