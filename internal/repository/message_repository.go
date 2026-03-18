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
		INSERT INTO messages (from_phone, type, body, media_id, created_at)
		VALUE ($1, $2, $3, $4, $5);
	`

	_, err := r.DB.Exec(
		query,
		msg.From,
		msg.Type,
		msg.Body,
		msg.MediaID,
		msg.Timestamp,
	)

	return err
}
