package database

type DatabasePort interface {
	StartDatabase() error
}
