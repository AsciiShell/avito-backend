package main

import (
	"net/http"

	"github.com/asciishell/avito-backend/internal/storage"
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

func (h Handler) Hello(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "not implemented", http.StatusNotFound)
}
