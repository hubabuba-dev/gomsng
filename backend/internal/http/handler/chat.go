package handler

import (
	"encoding/json"
	"msng/internal/service"
	"net/http"

	"github.com/google/uuid"
)

type GetChatsRequest struct {
	ChatType *string
	UserID   uuid.UUID
	Title    string
	Limit    int
}

type ChatHandler struct {
	chatService *service.ChatService
}

func NewChatHandler(chatService *service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) GetUserChats(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var request GetChatsRequest

	json_decoder := json.NewDecoder(r.Body)
	json_decoder.DisallowUnknownFields()
	err := json_decoder.Decode()
}
