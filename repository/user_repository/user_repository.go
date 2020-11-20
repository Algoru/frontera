package userrepository

import (
	"github.com/Algoru/frontera/domain/entity"
	"github.com/google/uuid"
)

// UserRepository
type UserRepository interface {
	CreateUser(*entity.User) (*entity.User, error)
	GetUser(uuid.UUID) (*entity.User, error)
	UpdateUser(uuid.UUID, *entity.User) (*entity.User, error)
	DeleteUser(uuid.UUID) (*entity.User, error)
	GetUserByEmail(string) (*entity.User, error)
}
