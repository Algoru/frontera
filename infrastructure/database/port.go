package database

import userrepository "github.com/Algoru/frontera/repository/user_repository"

// Port
type Port interface {
	StartDatabase() error
	userrepository.UserRepository
}
