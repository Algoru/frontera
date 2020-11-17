package userrepository

import (
	"github.com/Algoru/frontera/domain/entity"
	"github.com/google/uuid"
)

// UserRepository
type UserRepository interface {
	CreateUser(*User) (*entity.User, error)
	GetUser(uuid.UUID) (*entity.User, error)
	UpdateUser(uuid.UUID, User) (*entity.User, error)
	DeleteUser(uuid.UUID) (*entity.User, error)
}