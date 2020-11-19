package database

import (
	authrepository "github.com/Algoru/frontera/repository/auth_repository"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
)

// Port
type Port interface {
	StartDatabase() error
	userrepository.UserRepository
	authrepository.AuthRepository
}
