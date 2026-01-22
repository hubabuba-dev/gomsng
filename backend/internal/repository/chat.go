package repository

import (
	"context"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Conversation struct {
	ID        int64
	GUID      uuid.UUID
	Title     string
	Type      string
	DmUser1ID *uuid.UUID
	DmUser2ID *uuid.UUID
	CreatorID uuid.UUID
	CreatedAt time.Time
	DeletedAt time.Time
}

type Participants struct {
	ID             int64
	ConversationID int
	UserID         int
	CreatedAt      time.Time
}

type Message struct {
	ID             int64
	ConversationID int64
	SenderID       int64
	MessageBody    string
	CreatedAt      time.Time
}

type Chat struct {
	GUID  uuid.UUID
	Title string
	Type  string
}

type ChatRepository struct {
	pool *pgxpool.Pool
}

type ChatListFilter struct {
	Limit    int
	ChatType *string
}

type MessageFilter struct {
	Limit           int
	LastMessageTime time.Time
	LastId          int
}

func NewChatRepository(pool *pgxpool.Pool) *ChatRepository {
	return &ChatRepository{
		pool: pool,
	}
}

func (r *ChatRepository) GetAllChatList(ctx context.Context, user_id uuid.UUID, f ChatListFilter) (*[]Chat, error) {
	args := []any{user_id}

	if f.Limit <= 0 {
		f.Limit = 30
	} else if f.Limit >= 200 {
		f.Limit = 200
	}

	query := `
		SELECT conversations.guid, conversations.title, conversations.type
		FROM conversations
		JOIN participants ON participants.conversation_id = conversations.id
		JOIN users ON participants.user_id = users.id
		WHERE users.uuid = $1 and conversations.deleted_at IS NULL
	`

	if f.ChatType != nil {
		args = append(args, *f.ChatType)
		query = query + " AND conversations.type = $2 "
	}
	args = append(args, f.Limit)
	query = query + " LIMIT $" + strconv.Itoa(len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return &[]Chat{}, err
	}
	defer rows.Close()

	chats := make([]Chat, 0, f.Limit)

	for rows.Next() {
		var chat Chat
		if err := rows.Scan(&chat.GUID, &chat.Title, &chat.Type); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	return &chats, nil
}

func (r *ChatRepository) FindChat(ctx context.Context, chat_name string, chat_id uuid.UUID) (*Chat, error) {
	var chat Chat
	args := []any{}

	query := `
		SELECT conversations.guid, conversations.title, conversations.type
		FROM conversations
		JOIN participants ON participants.conversation_id = conversations.id
		JOIN users ON participants.user_id = users.id
		WHERE users.uuid = $1 AND conversations.guid = $2
		conversations.deleted_at IS NULL
	`

	if chat_name != "" {
		args = append(args, chat_name)
		query = query + " AND conversations.name = $2 "
	} else if chat_id != uuid.Nil {
		args = append(args, chat_id)
		query = query + " AND conversations.guid = $2 "
	}

	err := r.pool.QueryRow(ctx, query, args...).Scan(&chat.GUID, &chat.Title, &chat.Type)
	if err != nil {
		return &Chat{}, err
	}

	return &chat, nil
}

func (r *ChatRepository) OpenChat(ctx context.Context, guid uuid.UUID, user_id uuid.UUID) (*Conversation, error) {
	var chat Conversation
	var chat_type string
	var query_chat string
	query := `
		SELECT type
		FROM conversations
		WHERE guid = $1 
	`

	err := r.pool.QueryRow(ctx, query, guid).Scan(&chat_type)
	if err != nil {
		return &Conversation{}, err
	}

	switch chat_type {
	case "dm":
		query_chat = `
			SELECT id, guid, title, type, dm_user1_id, dm_user2_id, created_at, creator_id
			FROM conversations
			WHERE guid = $1 AND
			(dm_user1_id = $2 OR dm_user2_id = $2)
		`
	case "group":
		query_chat = `
			SELECT conversations.id, conversations.guid, conversations.title, conversations.type, conversations.creator_id, conversations.created_at
			FROM conversations
			JOIN participants ON conversations.id = participants.conversation_id 
			WHERE conversations.guid = $1
			AND participants.user_id = $2
			AND conversations.deleted_at IS NULL
		`
	}

	err = r.pool.QueryRow(ctx, query_chat, guid, user_id).Scan(&chat.ID, &chat.GUID, &chat.Title, &chat.Type, &chat.CreatedAt, &chat.DmUser1ID, &chat.DmUser2ID)
	if err != nil {
		return &Conversation{}, err
	}

	return &chat, nil
}

func (r *ChatRepository) GetMessages(ctx context.Context, conversation_id int, filter MessageFilter) (*[]Message, error) {
	args := []any{conversation_id}
	if filter.Limit <= 0 || filter.Limit >= 50 {
		filter.Limit = 15
	}

	query := `
		SELECT id, conversation_id, sender_id, message_body, created_at
		FROM messages
		WHERE conversation_id = $1
	`

	if filter.LastId != 0 && !filter.LastMessageTime.IsZero() {
		args = append(args, filter.LastId, filter.LastMessageTime)
		query += " AND (id, created_at) < ($2, $3)"
	}

	query += " ORDER BY created_at DESC, id DESC "

	args = append(args, filter.Limit)
	query += " LIMIT $" + strconv.Itoa(len(args))

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return &[]Message{}, err
	}
	defer rows.Close()

	messages := make([]Message, filter.Limit)
	for rows.Next() {
		var msg Message
		if err := rows.Scan(&msg.ID, &msg.ConversationID, &msg.SenderID, &msg.MessageBody, &msg.CreatedAt); err != nil {
			return &[]Message{}, err
		}
		messages = append(messages, msg)
	}

	return &messages, nil
}
