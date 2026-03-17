package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
)

type contactService struct {
	repo *repository.ContactRepository
}

func NewContactService(repo *repository.ContactRepository) *contactService {
	return &contactService{repo: repo}
}

func (s *contactService) SaveContact(phone string) {

	contact := models.Contact{
		Phone: phone,
		Name:  "",
	}
	s.repo.Save(contact)
}
