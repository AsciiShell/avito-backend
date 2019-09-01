package postgresqldb

import (
	"fmt"

	"github.com/asciishell/avito-backend/internal/chat"
	"github.com/asciishell/avito-backend/internal/message"
	"github.com/asciishell/avito-backend/internal/user"
	"github.com/asciishell/avito-backend/pkg/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Registry postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresStorage struct {
	DB *gorm.DB
}

type DBCredential struct {
	URL     string
	Debug   bool
	Migrate bool
}

func NewPostgresStorage(credential DBCredential) (*PostgresStorage, error) {
	var err error
	var db *gorm.DB
	logger := log.New()
	db, err = gorm.Open("postgres", credential.URL)
	if err != nil {
		return nil, errors.Wrapf(err, "can't connect to database, dsn %s", credential.URL)
	}
	if err = db.DB().Ping(); err != nil {
		return nil, errors.Wrapf(err, "can't ping database, dsn %s", credential.URL)
	}
	db.LogMode(credential.Debug)
	result := PostgresStorage{DB: db}
	if credential.Migrate {
		if err := result.Migrate(); err != nil {
			defer db.Close()
			return nil, errors.Wrapf(err, "can't apply migration")
		}
		logger.Info("Migration complete")
	}
	return &result, nil
}
func (p *PostgresStorage) constraintExists(table string, constraint string) bool {
	return p.DB.Exec(`SELECT 1 FROM pg_catalog.pg_constraint con
         INNER JOIN pg_catalog.pg_class rel ON rel.oid = con.conrelid
WHERE rel.relname = ? AND con.conname = ?;`, table, constraint).RowsAffected == 1
}

func (p *PostgresStorage) Migrate() error {
	t := p.DB.Begin()
	defer t.Commit()
	if err := p.DB.AutoMigrate(user.User{}, chat.Chat{}, message.Message{}).Error; err != nil {
		t.Rollback()
		return errors.Wrapf(err, "can't migrate")
	}
	type constrain struct {
		Table string
		Name  string
		Rule  string
	}
	constraints := []constrain{
		{Table: "user_chats", Name: "user_chats_chat_id_fkey", Rule: "FOREIGN KEY (chat_id) REFERENCES chats ON DELETE RESTRICT"},
		{Table: "user_chats", Name: "user_chats_user_id_fkey", Rule: "FOREIGN KEY (user_id) REFERENCES users ON DELETE RESTRICT"},
	}
	for _, v := range constraints {
		if !p.constraintExists(v.Table, v.Name) {
			if err := p.DB.Exec(fmt.Sprintf("ALTER TABLE %s ADD CONSTRAINT  %s %s", v.Table, v.Name, v.Rule)).Error; err != nil {
				t.Rollback()
				return errors.Wrapf(err, "can't apply constraint %s", v.Name)
			}
		}
	}
	return nil
}

func (p *PostgresStorage) CreateUser(u *user.User) error {
	if err := p.DB.Create(u).Error; err != nil {
		return errors.Wrapf(err, "can't create user %+v", u)
	}
	return nil
}

func (p *PostgresStorage) GetUser(u *user.User) error {
	if err := p.DB.Where(u).First(u).Error; err != nil {
		return errors.Wrapf(err, "can't get user %+v", u)
	}
	return nil
}

func (p *PostgresStorage) CreateChat(c *chat.Chat) error {
	if err := p.DB.Create(c).Error; err != nil {
		return errors.Wrapf(err, "can't create chat %+v", c)
	}
	return nil
}

func (p *PostgresStorage) CreateMessage(m *message.Message) error {
	if err := p.DB.Create(m).Error; err != nil {
		return errors.Wrapf(err, "can't create message %+v", m)
	}
	return nil
}

func (p *PostgresStorage) GetChatsFor(u user.User) ([]chat.Chat, error) {
	var result []chat.Chat
	if err := p.DB.Raw(`SELECT chats.id, chats.name, chats.created_at, (SELECT MAX(created_at) FROM messages WHERE messages.chat_id = chats.id) as last_message
FROM chats
         JOIN user_chats uc on chats.id = uc.chat_id
WHERE uc.user_id = ?
ORDER BY last_message DESC 
`, u.ID).Scan(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "can't get chats for user id %v", u.ID)
	}
	return result, nil
}

func (p *PostgresStorage) GetMessages(c chat.Chat) ([]message.Message, error) {
	var result []message.Message

	if err := p.DB.Raw(`SELECT *
FROM messages
WHERE chat_id = ?
ORDER BY created_at
`, c.ID).Scan(&result).Error; err != nil {
		return nil, errors.Wrapf(err, "can't read from row, chat id %v", c.ID)
	}
	return result, nil
}
