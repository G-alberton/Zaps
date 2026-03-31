package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/repository"
	"fmt"
	"log"
)

type ContactService struct {
	repo *repository.ContactRepository
}

func NewContactService(repo *repository.ContactRepository) *ContactService {
	return &ContactService{repo: repo}
}

func (s *ContactService) GetName(phone string) string {
	if s.repo != nil {
		contact, err := s.repo.GetByPhone(phone)
		if err == nil && contact != nil && contact.Name != "" {
			return contact.Name
		}
	}

	return phone
}

func (s *ContactService) SaveContact(phone, name string) error {

	if phone == "" {
		return fmt.Errorf("telefone vazio")
	}

	if len(phone) < 10 {
		return fmt.Errorf("telefone inválido")
	}

	if s.repo != nil {
		exists := s.repo.Exists(phone)
		if exists {
			log.Println("Contato já existe:", phone)
			return nil
		}
	}

	contact := models.Contact{
		Phone: phone,
		Name:  name,
	}

	if s.repo != nil {
		err := s.repo.Save(contact)
		if err != nil {
			log.Println("Erro ao salvar contato:", err)
			return err
		}
	}

	log.Println("Contato salvo:", phone)
	return nil
}
