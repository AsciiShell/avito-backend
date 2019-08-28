package storage

type Storage interface {
	Migrate(index int) error
}
