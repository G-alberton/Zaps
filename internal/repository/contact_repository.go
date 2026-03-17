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

func NewContatctRepository(db *sql.DB) *ContactRepository {
	return &ContactRepository{DB: db}
}

func (r *ContactRepository) Save(contact models.Contact) {
	query := `
		INSERT INTO contacts (phone, name)
		VALUES ($1, $2)
		ON CONFLICT (phone) DO NOTHING;
	`

	_, err := r.DB.Exec(query, contact.Phone, contact.Name)
	if err != nil {
		log.Println("Erro ao salvar contato:", err)
	}
}
