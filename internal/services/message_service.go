package services

import (
	"ZAPS/internal/models"
	"ZAPS/internal/pagination"
	"ZAPS/internal/repository"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type MessageService struct {
	repo     *repository.MessageRepository
	messages []models.Message
	mu       sync.Mutex
}

func NewMessageService(repo *repository.MessageRepository) *MessageService {
	return &MessageService{
		repo:     repo,
		messages: []models.Message{},
	}
}

func (s *MessageService) SaveMessage(msg models.Message) error {
	if msg.From == "" {
		return nil
	}

	switch msg.Type {
	case "text":
		if msg.Body == "" {
			return fmt.Errorf("Mensagem de Texto Vazio")
		}
	case "image", "audio", "document":
		if msg.MediaID == "" {
			log.Println("Mídia sem ID")
			return nil
		}
	default:
		log.Println("Tipo não suportado:", msg.Type)
		return nil
	}

	if msg.Timestamp.IsZero() {
		msg.Timestamp = time.Now()
	}

	s.mu.Lock()
	s.messages = append(s.messages, msg)
	s.mu.Unlock()

	if s.repo != nil {
		err := s.repo.Save(msg)
		if err != nil {
			log.Println("Erro ao salvar mensagem:", err)
		}
	}

	log.Printf("Mensagem salva | From: %s | Type: %s | Direction: %s", msg.From, msg.Type, msg.Direction)
	return nil
}

func (s *MessageService) GetByConversation(conversationID string) []models.Message {
	s.mu.Lock()
	messagesCopy := make([]models.Message, len(s.messages))
	copy(messagesCopy, s.messages)
	s.mu.Unlock()

	var result []models.Message

	for _, msg := range messagesCopy {
		if msg.ConversationID == conversationID {
			result = append(result, msg)
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Timestamp.Before(result[j].Timestamp)
	})

	return result
}

func (s *MessageService) GetLastMessage(conversationID string) *models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	var last *models.Message

	for i := range s.messages {
		msg := &s.messages[i]

		if msg.ConversationID == conversationID {
			if last == nil || msg.Timestamp.After(last.Timestamp) {
				last = msg
			}
		}
	}

	return last
}

func (s *MessageService) GetAll() []models.Message {
	s.mu.Lock()
	defer s.mu.Unlock()

	copySlice := make([]models.Message, len(s.messages))
	copy(copySlice, s.messages)

	return copySlice
}

func (s *MessageService) CountUnread(conversationID string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	count := 0

	for _, msg := range s.messages {
		if msg.ConversationID == conversationID && msg.Direction == "inbound" {
			count++
		}
	}

	return count
}

func (s *MessageService) MarkAsRead(conversationID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i := range s.messages {
		if s.messages[i].ConversationID == conversationID {
			s.messages[i].Direction = "read"
		}
	}
}

func (s *MessageService) ListMessagesPaginated(p pagination.Pagination) (pagination.Response[models.Message], error) {

	p.Normalize()

	messages, err := s.repo.ListPaginated(p.Cursor, p.Limit+1)
	if err != nil {
		return pagination.Response[models.Message]{}, fmt.Errorf("repository not initialized")
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
