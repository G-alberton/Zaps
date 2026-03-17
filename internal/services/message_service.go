package services

import (
	"log"

	"ZAPS/internal/models"
	"ZAPS/internal/repository"
)

type MessageService struct {
	repo *repository.MessageRepository
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) SaveMessage(from, msgType, body string) error {
	if from == "" {
		return nil
	}

	message := models.Message{
		From: from,
		Type: msgType,
		Body: body,
	}

	err := s.repo.Save(message)
	if err != nil {
		log.Println("Erro ao Salvar mensagem:", err)
		return err
	}

	log.Println("Mensagem Salva:", from, body)
	return nil
}
