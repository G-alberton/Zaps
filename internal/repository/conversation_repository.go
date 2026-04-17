package repository

import (
	"ZAPS/internal/models"
	"database/sql"
)

type ConversationRepository struct {
	DB *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{DB: db}
}

func (r *ConversationRepository) GetByContact(contact string) (*models.Conversation, error) {
	query := `
		SELECT id, contact_phone
		FROM conversations
		WHERE contact_phone = $1
		LIMIT 1
	`

	var c models.Conversation

	err := r.DB.QueryRow(query, contact).Scan(
		&c.ID,
		&c.Contact,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *ConversationRepository) Create(contact string) (string, error) {
	query := `
		INSERT INTO conversations (contact_phone)
		Values ($1)
		RETURNING id
	`

	var id string
	err := r.DB.QueryRow(query, contact).Scan(&id)

	return id, err
}

func (r *ConversationRepository) UpdateLastMessage(conversationID string) error {
	query := `
		UPDATE conversations
		SET last_message_at = now()
		where id = $1
	`

	_, err := r.DB.Exec(query, conversationID)
	return err
}

func (r *ConversationRepository) ListAll() ([]models.Conversation, error) {
	rows, err := r.DB.Query(`
		SELECT id, contact_phone, status, last_message_at, created_at
		FROM conversation
		ORDER BY last_message_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Conversation

	for rows.Next() {
		var c models.Conversation

		err := rows.Scan(
			&c.ID,
			&c.Contact,
			&c.Status,
			&c.LastMessageAt,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		list = append(list, c)
	}

	return list, nil
}
