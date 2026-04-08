package repository

import (
	"ZAPS/internal/models"
	"database/sql"
	"time"
)

type MessageRepository struct {
	DB *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

func (r *MessageRepository) Save(msg models.Message) error {
	query := `
		INSERT INTO messages (from_phone, type, body, media_id, created_at, conversation_id)
		VALUES ($1, $2, $3, $4, $5, $6);
	`

	_, err := r.DB.Exec(
		query,
		msg.From,
		msg.Type,
		msg.Body,
		msg.MediaID,
		msg.Timestamp,
		msg.ConversationID,
	)

	return err
}

func (r *MessageRepository) ListPaginated(cursor time.Time, limit int) ([]models.Message, error) {

	rows, err := r.DB.Query(`
		SELECT id, from_phone, type, body, media_id, created_at, conversation_id
		FROM messages
		WHERE created_at > $1
		ORDER BY created_at ASC
		LIMIT $2
	`, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.From,
			&m.Type,
			&m.Body,
			&m.MediaID,
			&m.Timestamp,
			&m.ConversationID,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *MessageRepository) ListPaginatedByConversation(
	conversationID string,
	cursor time.Time,
	limit int,
) ([]models.Message, error) {

	rows, err := r.DB.Query(`
		SELECT id, from_phone, type, body, media_id, created_at, conversation_id
		FROM messages
		WHERE conversation_id = $1 AND id > $2
		ORDER BY created_at ASC
		LIMIT $3
	`, conversationID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []models.Message

	for rows.Next() {
		var m models.Message
		if err := rows.Scan(
			&m.ID,
			&m.From,
			&m.Type,
			&m.Body,
			&m.MediaID,
			&m.Timestamp,
			&m.ConversationID,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}
