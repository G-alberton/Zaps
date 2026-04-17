package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
)

type ConversationService struct {
	repo *repository.ConversationRepository
}

func NewConversationService(repo *repository.ConversationRepository) *ConversationService {
	return &ConversationService{
		repo: repo,
	}
}

func (s *ConversationService) GetOrCreate(contact string) (string, error) {
	conv, err := s.repo.GetByContact(contact)
	if err != nil {
		return "", err
	}

	if conv != nil {
		return conv.ID, nil
	}

	return s.repo.Create(contact)
}

func (s *ConversationService) List() ([]models.Conversation, error) {
	return s.repo.ListAll()
}
