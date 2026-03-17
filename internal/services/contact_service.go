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

func (s *ContactService) SaveContact(phone, name string) error {

	if phone == "" {
		return fmt.Errorf("telefone vazio")
	}

	if len(phone) < 10 {
		return fmt.Errorf("telefone inválido")
	}

	exists := s.repo.Exists(phone)
	if exists {
		log.Println("Contato já existe:", phone)
		return nil
	}

	contact := models.Contact{
		Phone: phone,
		Name:  name,
	}

	err := s.repo.Save(contact)
	if err != nil {
		log.Println("Erro ao salvar contato:", err)
		return err
	}

	log.Println("Contato salvo:", phone)
	return nil
}
