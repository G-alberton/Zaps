package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
)

type ContactService struct {
	repo *repository.ContactRepository
}

func NewContactService(repo *repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) SaveContact(phone string) {

	contact := models.Contact{
		Phone: phone,
		Name:  "",
	}
	s.repo.Save(contact)
}
