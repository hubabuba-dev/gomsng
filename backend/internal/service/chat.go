package service

import (
	"context"
	"log"
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

func (s *ChatService) GetUserChats(ctx context.Context, chats Chats) (*[]repository.Chat, error) {
	if chats.ChatType != nil {
		if *chats.ChatType != "group" && *chats.ChatType != "dm" {
			return nil, ErrWrongChatType
		}
	}

	filter := repository.ChatListFilter{
		Limit:    chats.Limit,
		ChatType: chats.ChatType,
	}

	chat_list, err := s.chat_repo.GetAllChatList(ctx, chats.UserID, filter)
	if err != nil {
		return nil, err
	}

	return chat_list, nil
}

func (s *ChatService) OpenUserChat(ctx context.Context, guid uuid.UUID, user_id uuid.UUID) {
	chat, err := s.chat_repo.OpenChat(ctx, guid, user_id)
	if err != nil {
		//return nil, err
	}
	log.Println(chat)
	//return chat, err
}

func (s *ChatService) GetMessages(ctx context.Context) {

}
