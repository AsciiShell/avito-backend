package storage

import (
	"github.com/asciishell/avito-backend/internal/chat"
	"github.com/asciishell/avito-backend/internal/message"
	"github.com/asciishell/avito-backend/internal/user"
)

type Storage interface {
	Migrate() error
	CreateUser(u *user.User) error
	GetUser(u *user.User) error
	CreateChat(c *chat.Chat) error
	CreateMessage(m *message.Message) error
	GetChatsFor(u user.User) ([]chat.Chat, error)
	GetMessages(c chat.Chat) ([]message.Message, error)
}
