package chat

import (
	"fmt"
	"time"

	"github.com/asciishell/avito-backend/internal/user"
)

type Chat struct {
	ID        uint        `json:"id" gorm:"primary_key"`
	Name      string      `json:"name" gorm:"NOT NULL;unique"`
	Users     []user.User `json:"users" gorm:"many2many:user_chats;"`
	CreatedAt time.Time   `json:"created_at" gorm:"NOT NULL;DEFAULT:now()"`
}

func (c Chat) ShortJSON() []byte {
	return []byte(fmt.Sprintf("{\"id\": %v}", c.ID))
}
