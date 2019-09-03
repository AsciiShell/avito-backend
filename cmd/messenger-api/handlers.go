package main

import (
	"encoding/json"
	"io"
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

func (h Handler) Write(w io.Writer, b []byte) {
	if _, err := w.Write(b); err != nil {
		h.logger.Errorf("can't write data: %+v", err)
	}
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

	h.Write(w, userData.ShortJSON())
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

	h.Write(w, chatData.ShortJSON())
}
func (h Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var messageData message.Message
	if err := json.NewDecoder(r.Body).Decode(&messageData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.storage.CreateMessage(&messageData); err != nil {
		h.logger.Errorf("can't create message: %+v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	h.Write(w, messageData.ShortJSON())
}

//nolint:dupl
func (h Handler) GetChats(w http.ResponseWriter, r *http.Request) {
	var userInfo user.CreationUser
	if err := json.NewDecoder(r.Body).Decode(&userInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userData := userInfo.Convert()
	chats, err := h.storage.GetChatsFor(userData)
	switch err {
	case nil:
	case postgresqldb.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		h.logger.Errorf("can't get chats for %v: %+v", userData, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(chats)
	if err != nil {
		h.logger.Errorf("can't write message info: %+v", err)
		return
	}
}

//nolint:dupl
func (h Handler) GetMessages(w http.ResponseWriter, r *http.Request) {
	var chatInfo chat.CreationChat
	if err := json.NewDecoder(r.Body).Decode(&chatInfo); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	chatData := chatInfo.Convert()
	messages, err := h.storage.GetMessages(chatData)
	switch err {
	case nil:
	case postgresqldb.ErrNotFound:
		w.WriteHeader(http.StatusNotFound)
		return
	default:
		h.logger.Errorf("can't get messages for %v: %+v", chatData, err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		h.logger.Errorf("can't write message info: %+v", err)
		return
	}
}
