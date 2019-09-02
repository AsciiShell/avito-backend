package main

import (
	"encoding/json"
	"net/http"

	"github.com/asciishell/avito-backend/internal/message"

	"github.com/asciishell/avito-backend/internal/chat"
	"github.com/asciishell/avito-backend/internal/postgresqldb"
	"github.com/asciishell/avito-backend/internal/storage"
	"github.com/asciishell/avito-backend/internal/user"
	"github.com/asciishell/avito-backend/pkg/log"
)

type Handler struct {
	storage storage.Storage
	logger  log.Logger
}

func NewHandler(l log.Logger, s storage.Storage) *Handler {
	h := Handler{storage: s, logger: l}
	return &h
}

func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var userData user.User
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.storage.CreateUser(&userData); err != nil {
		h.logger.Errorf("can't create user: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(userData.ShortJSON()); err != nil {
		h.logger.Errorf("can't write user info: %+v", err)
		return
	}
}
func (h Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var chatInfo chat.CreationChat
	if err := json.NewDecoder(r.Body).Decode(&chatInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatData := chatInfo.Convert()
	if err := h.storage.CreateChat(&chatData); err != nil {
		h.logger.Errorf("can't create chatData: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if _, err := w.Write(chatData.ShortJSON()); err != nil {
		h.logger.Errorf("can't write chatData info: %+v", err)
		return
	}
}
func (h Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var messageInfo message.CreationMessage
	if err := json.NewDecoder(r.Body).Decode(&messageInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	messageData := messageInfo.Convert()
	if err := h.storage.CreateMessage(&messageData); err != nil {
		h.logger.Errorf("can't create message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(messageData.ShortJSON()); err != nil {
		h.logger.Errorf("can't write message info: %+v", err)
		return
	}
}
func (h Handler) GetChats(w http.ResponseWriter, r *http.Request) {
	var userInfo user.CreationUser
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userData := userInfo.Convert()
	chats, err := h.storage.GetChatsFor(userData)
	if err != nil {
		switch err {
		case postgresqldb.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			h.logger.Errorf("can't get chats for %v: %+v", userData, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		h.logger.Errorf("can't write message info: %+v", err)
		return
	}
}
func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var chatData chat.Chat
	if err := json.NewDecoder(r.Body).Decode(&chatData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	messages, err := h.storage.GetMessages(chatData)
	if err != nil {
		switch err {
		case postgresqldb.ErrNotFound:
			w.WriteHeader(http.StatusNotFound)
		default:
			h.logger.Errorf("can't get messages for %v: %+v", chatData, err)
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		h.logger.Errorf("can't write message info: %+v", err)
		return
	}
}
