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
	DmUser1ID uuid.UUID
	DmUser2ID uuid.UUID
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
