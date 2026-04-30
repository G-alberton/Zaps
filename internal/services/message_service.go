package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/pagination"
	"ZAPS/internal/repository"
	"ZAPS/internal/websocket"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"
)

type MessageService struct {
	repo             *repository.MessageRepository
	conversationRepo *repository.ConversationRepository
	hub              *websocket.Hub
	DB               *sql.DB
}

func NewMessageService(
	repo *repository.MessageRepository,
	conversationRepo *repository.ConversationRepository,
	hub *websocket.Hub,
) *MessageService {
	return &MessageService{
		repo:             repo,
		conversationRepo: conversationRepo,
		hub:              hub,
	}
}

func (s *MessageService) SaveMessage(msg models.Message) error {
	if msg.From == "" {
		return fmt.Errorf("from vazio")
	}

	switch msg.Type {
	case "text":
		if msg.Body == "" {
			return fmt.Errorf("mensagem vazia")
		}
	case "image", "audio", "document":
		if msg.MediaID == "" {
			return fmt.Errorf("media sem ID")
		}
	default:
		return fmt.Errorf("tipo não suportado")
	}

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	if msg.Status == "" {
		msg.Status = "sent"
	}

	if msg.Direction == "" {
		msg.Direction = "outbound"
	}

	if err := s.repo.Save(msg); err != nil {
		return err
	}

	if s.conversationRepo != nil {
		_ = s.conversationRepo.UpdateLastMessage(msg.ConversationID)
	}

	if s.hub != nil {
		data, _ := json.Marshal(msg)

		s.hub.Broadcast <- websocket.MessagePayload{
			ConversationID: msg.ConversationID,
			Data:           data,
		}
	}

	log.Printf("Mensagem salva | %s | %s", msg.From, msg.Type)

	return nil
}

func (s *MessageService) GetByConversation(conversationID string) ([]models.Message, error) {
	return s.repo.ListByConversation(conversationID)
}

func (s *MessageService) ListMessagesPaginated(
	p pagination.Pagination,
) (pagination.Response[models.Message], error) {

	p.Normalize()

	messages, err := s.repo.ListPaginated(p.Cursor, p.Limit+1)
	if err != nil {
		return pagination.Response[models.Message]{}, err
	}

	hasMore := false
	if len(messages) > p.Limit {
		hasMore = true
		messages = messages[:p.Limit]
	}

	var nextCursor *time.Time
	if len(messages) > 0 {
		t := messages[len(messages)-1].Timestamp
		nextCursor = &t
	}

	return pagination.Response[models.Message]{
		Data:       messages,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *MessageService) ListMessagesByConversationPaginated(
	conversationID string,
	p pagination.Pagination,
) (pagination.Response[models.Message], error) {

	p.Normalize()

	messages, err := s.repo.ListPaginatedByConversation(
		conversationID,
		p.Cursor,
		p.Limit+1,
	)
	if err != nil {
		return pagination.Response[models.Message]{}, err
	}

	hasMore := false
	if len(messages) > p.Limit {
		hasMore = true
		messages = messages[:p.Limit]
	}

	var nextCursor *time.Time
	if len(messages) > 0 {
		t := messages[len(messages)-1].Timestamp
		nextCursor = &t
	}

	return pagination.Response[models.Message]{
		Data:       messages,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

func (s *MessageService) GetLastMessage(conversationID string) (*models.Message, error) {
	return s.repo.GetLastMessage(conversationID)
}

func (s *MessageService) CountUnread(conversationID string) (int, error) {
	return s.repo.CountUnread(conversationID)
}

func (s *MessageService) MarkAsRead(conversationID string) error {
	return s.repo.MarkAsRead(conversationID)
}

func (s *MessageService) UpdateStatus(messageID string, status string) error {
	return s.repo.UpdateStatus(messageID, status)
}

func (s *MessageService) Exists(messageID string) (bool, error) {
	return s.repo.Exists(messageID)
}

func (s *MessageService) GetByID(id string) (models.Message, error) {
	var msg models.Message

	query := `
        SELECT id, conversation_id, "from", type, body, media_url, direction, status, timestamp
        FROM messages
        WHERE id = $1
    `

	err := s.DB.QueryRow(query, id).Scan(
		&msg.ID,
		&msg.ConversationID,
		&msg.From,
		&msg.Type,
		&msg.Body,
		&msg.MediaURL,
		&msg.Direction,
		&msg.Status,
		&msg.Timestamp,
	)

	if err != nil {
		return models.Message{}, err
	}

	return msg, nil
}
