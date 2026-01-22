package handler

import (
	"encoding/json"
	"msng/internal/service"
	"net/http"

	"github.com/google/uuid"
)

type ChatHandler struct {
	chatService *service.ChatService
}

type GetChatsRequest struct {
	GUID     uuid.UUID
	Limit    int
	ChatType *string
}

type FindChatRequest struct {
	GUID     uuid.UUID
	ChatName string
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
	err := json_decoder.Decode(&request)
	if err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}

	chats := service.Chats{
		Limit:    request.Limit,
		ChatType: request.ChatType,
	}
	chats_list, err := h.chatService.GetUserChats(ctx, chats)
	if err != nil {
		http.Error(w, "", http.StatusBadGateway)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(chats_list)
}

///func (h *ChatHandler) FindChat(w http.ResponseWriter, r *http.Request) {
///	ctx := r.Context()
///	var request FindChatRequest
///
///	json_decoder := json.NewDecoder(r.Body)
///	json_decoder.DisallowUnknownFields()
///	err := json_decoder.Decode(&request)
///	if err != nil {
///		http.Error(w, "invalid json", http.StatusBadRequest)
///		return
///	}
///}
