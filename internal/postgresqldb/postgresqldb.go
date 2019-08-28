package postgresqldb

import (
	"github.com/asciishell/avito-backend/pkg/log"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// Registry postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type PostgresStorage struct {
	DB *gorm.DB
}

func (p *PostgresStorage) Migrate() error {
	panic("implement me")
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
