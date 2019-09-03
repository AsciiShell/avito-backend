package user

import (
	"fmt"
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"PRIMARY_KEY"`
	Username  string    `json:"username" gorm:"NOT NULL;UNIQUE"`
	CreatedAt time.Time `json:"created_at" gorm:"NOT NULL;DEFAULT:now()"`
}

type CreationUser struct {
	User uint `json:"user"`
}

func (u User) ShortJSON() []byte {
	return []byte(fmt.Sprintf("{\"id\": %v}", u.ID))
}

func (c CreationUser) Convert() User {
	return User{ID: c.User}
}
