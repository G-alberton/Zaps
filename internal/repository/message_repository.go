package repository

import (
	"ZAPS/internal/models"
	"database/sql"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) Save(msg models.Message) error {
	query := `
		INSERT INTO messages (phone, type, body)
		VALUE ($1, $2, $3);
	`

	_, err := r.DB.Exec(query, msg.From, msg.Type, msg.Body)
	if err != nil {
		return err
	}

	return nil
}
