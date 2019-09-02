package message

import (
	"fmt"
	"time"

	"github.com/asciishell/avito-backend/internal/chat"
	"github.com/asciishell/avito-backend/internal/user"
)

type Message struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Chat      chat.Chat `json:"chat" gorm:"NOT NULL;foreignkey:Chat"`
	ChatID    uint      `json:"-" gorm:"NOT NULL" sql:"type:integer REFERENCES chats(id)"`
	Author    user.User `json:"author" gorm:"foreignkey:Author"`
	AuthorID  uint      `json:"-" gorm:"NOT NULL" sql:"type:integer REFERENCES users(id)"`
	Text      string    `json:"text" gorm:"NOT NULL"`
	CreatedAT time.Time `json:"created_at" gorm:"NOT NULL;DEFAULT:now()"`
}

func (m Message) ShortJSON() []byte {
	return []byte(fmt.Sprintf("{\"id\": %v}", m.ID))
}

type CreationMessage struct {
	Chat   uint   `json:"chat"`
	Author uint   `json:"author"`
	Text   string `json:"text"`
}

func (c CreationMessage) Convert() Message {
	return Message{
		ChatID:   c.Chat,
		AuthorID: c.Author,
		Text:     c.Text,
	}
}
