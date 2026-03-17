package repository

import (
	"ZAPS/internal/models"
	"database/sql"
	"log"
)

// aqui a gente se comunica com o banco de dados
type ContactRepository struct {
	DB *sql.DB
}

func NewContactRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{DB: db}
}

func (r *ContactRepository) Save(contact models.Contact) error {
	query := `
		INSERT INTO contacts (phone, name)
		VALUES ($1, $2)
		ON CONFLICT (phone) DO NOTHING;
	`

	_, err := r.DB.Exec(query, contact.Phone, contact.Name)
	if err != nil {
		return err
	}

	return nil
}

func (r *ContactRepository) Exists(phone string) bool {
	query := `SELECT EXISTS(SELECT 1 FROM contacts WHERE phone = $1)`

	var exists bool

	err := r.DB.QueryRow(query, phone).Scan(&exists)
	if err != nil {
		log.Println("Erro ao verificar contato:", err)
		return false
	}

	return exists
}
