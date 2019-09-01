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
	panic("implement me")
}

func (p *PostgresStorage) GetUser(u *user.User) error {
	panic("implement me")
}

func (p *PostgresStorage) CreateChat(c *chat.Chat) error {
	panic("implement me")
}

func (p *PostgresStorage) CreateMessage(m *message.Message) error {
	panic("implement me")
}

func (p *PostgresStorage) GetChatsFor(u user.User) ([]chat.Chat, error) {
	panic("implement me")
}

func (p *PostgresStorage) GetMessages(c chat.Chat) ([]message.Message, error) {
	panic("implement me")
}
