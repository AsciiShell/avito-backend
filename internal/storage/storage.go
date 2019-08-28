package storage

type Storage interface {
	Migrate() error
}
