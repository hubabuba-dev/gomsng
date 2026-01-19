package service

import (
	"context"
	"msng/internal/repository"

	"github.com/google/uuid"
)

type Chats struct {
	UserID   uuid.UUID
	Limit    int
	ChatType *string
}

type ChatService struct {
	chat_repo *repository.ChatRepository
}

func NewChatRepository(repo *repository.ChatRepository) *ChatService {
	return &ChatService{
		chat_repo: repo,
	}
}

func (s *ChatService) GetUserChats(ctx context.Context, chats Chats) (*[]Chats, error) {
	if chats.ChatType != nil {
		if *chats.ChatType != "group" && *chats.ChatType != "dm" {
			return nil, ErrWrongChatType
		}
	}

	filter := repository.ChatListFilter{
		Limit:    chats.Limit,
		ChatType: chats.ChatType,
	}

	s.chat_repo.GetAllChatList(chats.UserID, filter)

	return nil, nil
}
