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
		INSERT INTO messages (
		id,
		conversation_id,
		from_phone, 
		type, 
		body, 
		media_id,
		media_url,
		sent_at,
		direction,
		status,
		read 
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	_, err := r.DB.Exec(
		query,
		msg.ID,
		msg.ConversationID,
		msg.From,
		msg.Type,
		msg.Body,
		msg.MediaID,
		msg.MediaURL,
		msg.Timestamp,
		msg.Direction,
		msg.Status,
		msg.Read,
	)

	return err
}

func (r *MessageRepository) ListPaginated(cursor time.Time, limit int) ([]models.Message, error) {

	rows, err := r.DB.Query(`
		SELECT id, from_phone, type, body, media_id, media_url, sent_at, direction, status, read, conversation_id
		FROM messages
		WHERE sent_at > $1
		ORDER BY sent_at ASC
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
			&m.MediaURL,
			&m.Timestamp,
			&m.Direction,
			&m.Status,
			&m.Read,
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
		SELECT id, from_phone, type, body, media_id, media_url, sent_at, direction, status, read, conversation_id
		FROM messages
		WHERE conversation_id = $1 AND sent_at > $2
		ORDER BY sent_at ASC
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
			&m.MediaURL,
			&m.Timestamp,
			&m.Direction,
			&m.Status,
			&m.Read,
			&m.ConversationID,
		); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}

	return messages, nil
}

func (r *MessageRepository) ListByConversation(conversationID string) ([]models.Message, error) {
	rows, err := r.DB.Query(`
		SELECT id, from_phone, type, body, media_id, media_url, sent_at, direction, status, read, conversation_id
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at ASC
	`, conversationID)
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
			&m.MediaURL,
			&m.Timestamp,
			&m.Direction,
			&m.Status,
			&m.Read,
			&m.ConversationID,
		); err != nil {
			return nil, err
		}

		messages = append(messages, m)
	}

	return messages, nil
}

func (r *MessageRepository) GetLastMessage(conversationID string) (*models.Message, error) {
	row := r.DB.QueryRow(`
		SELECT id, conversation_id, from_phone, type, body, media_id, media_url, sent_at, direction, status, read
		FROM messages
		WHERE conversation_id = $1
		ORDER BY sent_at DESC
		LIMIT 1
	`, conversationID)

	var m models.Message

	err := row.Scan(
		&m.ID,
		&m.ConversationID,
		&m.From,
		&m.Type,
		&m.Body,
		&m.MediaID,
		&m.MediaURL,
		&m.Timestamp,
		&m.Direction,
		&m.Status,
		&m.Read,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &m, nil
}

func (r *MessageRepository) CountUnread(conversationID string) (int, error) {
	var count int

	err := r.DB.QueryRow(`
		SELECT COUNT(*)
		FROM messages
		WHERE conversation_id = $1
		AND direction = 'inbound'
		AND read = false
	`, conversationID).Scan(&count)

	return count, err
}

func (r *MessageRepository) MarkAsRead(conversationID string) error {
	_, err := r.DB.Exec(`
		UPDATE messages
		SET read = true
		WHERE conversation_id = $1
		AND direction = 'inbound'
		AND read = false
	`, conversationID)

	return err
}

func (r *MessageRepository) UpdateStatus(id string, status string) error {
	query := `UPDATE messages SET status = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, status, id)
	return err
}

func (s *MessageRepository) Exists(messageID string) (bool, error) {
	var count int

	err := s.DB.QueryRow(`
		SELECT COUNT(1)
		FROM messages
		WHERE id = ?
	`, messageID).Scan(&count)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
