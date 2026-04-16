package repository

import "database/sql"

type ConversationRepository struct {
	DB *sql.DB
}

func NewConversationRepository(db *sql.DB) *ConversationRepository {
	return &ConversationRepository{DB: db}
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
