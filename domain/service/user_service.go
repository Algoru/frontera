package service

import (
	"time"

	"github.com/Algoru/frontera/domain/entity"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
	"github.com/google/uuid"
)

// UserService
type UserService interface {
	CreateUser(*userrepository.User) (*entity.User, error)
	GetUser(uuid.UUID) (*entity.User, error)
	UpdateUser(uuid.UUID, *entity.User) (*entity.User, error)
	DeleteUser(uuid.UUID) (*entity.User, error)
}

type service struct {
	repo userrepository.UserRepository
}

func NewUserService(r userrepository.UserRepository) UserService {
	return &service{repo: r}
}

func (s *service) CreateUser(u *userrepository.User) (*entity.User, error) {
	prepared, err := u.Prepare()
	if err != nil {
		return nil, err
	}

	prepared.CreatedAt = time.Now()
	prepared.UpdatedAt = prepared.CreatedAt

	return s.repo.CreateUser(prepared)
}

func (s *service) GetUser(userID uuid.UUID) (*entity.User, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	user.RemoveSensible()

	return user, nil
}

func (s *service) UpdateUser(userID uuid.UUID, update *entity.User) (*entity.User, error) {
	return nil, nil
}

func (s *service) DeleteUser(userID uuid.UUID) (*entity.User, error) {
	return s.repo.DeleteUser(userID)
}
