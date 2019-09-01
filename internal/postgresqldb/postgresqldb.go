package postgresqldb

import (
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

func (p *PostgresStorage) Migrate() error {
	if err := p.DB.AutoMigrate(user.User{}, chat.Chat{}, message.Message{}).Error; err != nil {
		return errors.Wrapf(err, "can't migrate")
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
	rows, err := p.DB.Raw(`SELECT chats.id, chats.name, chats.created_at, (SELECT MAX(created_at) FROM messages WHERE messages.chat_id = chats.id) as last_message
FROM chats
         JOIN user_chats uc on chats.id = uc.chat_id
WHERE uc.user_id = ?
ORDER BY last_message DESC 
`, u.ID).Rows()
	if err != nil {
		return nil, errors.Wrapf(err, "can't get chats for user id %v", u.ID)
	}
	defer rows.Close()
	var result []chat.Chat
	for rows.Next() {
		var ch chat.Chat
		if err := p.DB.ScanRows(rows, &ch); err != nil {
			return nil, errors.Wrapf(err, "can't read from row, client id %v", u.ID)
		}
		result = append(result, ch)
	}
	return result, nil
}

func (p *PostgresStorage) GetMessages(c chat.Chat) ([]message.Message, error) {
	rows, err := p.DB.Raw(`SELECT *
FROM messages
WHERE chat_id = ?
ORDER BY created_at
`, c.ID).Rows()
	if err != nil {
		return nil, errors.Wrapf(err, "can't get messages for chat id %v", c.ID)
	}
	defer rows.Close()
	var result []message.Message
	for rows.Next() {
		var m message.Message
		if err := p.DB.ScanRows(rows, &m); err != nil {
			return nil, errors.Wrapf(err, "can't read from row, chat id %v", c.ID)
		}
		result = append(result, m)
	}
	return result, nil
}
