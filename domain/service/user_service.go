package service

import (
	"fmt"
	"time"

	"github.com/Algoru/frontera/configuration"

	"github.com/Algoru/frontera/domain/entity"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
	"github.com/badoux/checkmail"
	"github.com/google/uuid"
)

// UserService
type UserService interface {
	CreateUser(*userrepository.User) (*entity.User, error)
	GetUser(uuid.UUID) (*entity.User, error)
	UpdateUser(uuid.UUID, *entity.User) (*entity.User, error)
	DeleteUser(uuid.UUID) (*entity.User, error)
	HasRequiredFields(*userrepository.User) []string
	GetUserByEmail(string) (*entity.User, error)
}

type userService struct {
	repo userrepository.UserRepository
}

func NewUserService(r userrepository.UserRepository) UserService {
	return &userService{repo: r}
}

func (s *userService) CreateUser(u *userrepository.User) (*entity.User, error) {
	prepared, err := u.Prepare()
	if err != nil {
		return nil, err
	}

	prepared.CreatedAt = time.Now()
	prepared.UpdatedAt = prepared.CreatedAt

	return s.repo.CreateUser(prepared)
}

func (s *userService) GetUser(userID uuid.UUID) (*entity.User, error) {
	user, err := s.repo.GetUser(userID)
	if err != nil {
		return nil, err
	}
	user.RemoveSensible()

	return user, nil
}

func (s *userService) UpdateUser(userID uuid.UUID, update *entity.User) (*entity.User, error) {
	return nil, nil
}

func (s *userService) DeleteUser(userID uuid.UUID) (*entity.User, error) {
	return s.repo.DeleteUser(userID)
}

func (s *userService) HasRequiredFields(u *userrepository.User) []string {
	errors := make([]string, 0)

	if err := checkmail.ValidateFormat(u.Email); err != nil {
		err := "invalid email: " + err.Error()
		errors = append(errors, err)
	}

	minPasswordLength := int(configuration.GetConfiguration().Security.MinPasswordLength)
	if len(u.Password) < minPasswordLength {
		err := fmt.Sprintf("the password must contain at least %d characters", minPasswordLength)
		errors = append(errors, err)
	}

	userConfig := configuration.GetConfiguration().User
	if !userConfig.AllowNullPayload && u.Payload == nil {
		errors = append(errors, "payload cannot be null")
	}

	// TODO (@Algoru): Add configuration defined users vaidation for payload fields

	if len(errors) == 0 {
		return nil
	}

	return errors
}

func (s *userService) GetUserByEmail(email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(email)
}
