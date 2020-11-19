package userrepository

import (
	"strings"
	"time"

	"github.com/Algoru/frontera/configuration"

	"github.com/Algoru/frontera/domain/entity"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const encryptPasswordCost = 12

// User
type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Payload   bson.M             `bson:"payload"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// ToRawEntity
func (u *User) ToRawEntity() entity.User {
	// TODO (@Algoru): Should check this error ?
	userID, _ := uuid.Parse(u.UserID)

	return entity.User{
		ID:        u.ID,
		UserID:    userID,
		Email:     u.Email,
		Password:  u.Password,
		Payload:   u.Payload,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Sanitize
func (u *User) Sanitize() {
	u.Email = strings.ToLower(u.Email)
	u.Email = bluemonday.StrictPolicy().Sanitize(u.Email)
}

// EncryptPassword
func (u *User) EncryptPassword() (string, error) {
	passwordBytes := []byte(u.Password)

	cost := configuration.GetConfiguration().Security.BCryptCost
	hash, err := bcrypt.GenerateFromPassword(passwordBytes, int(cost))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

// Prepare
func (u *User) Prepare() (*User, error) {
	userID, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	password, err := u.EncryptPassword()
	if err != nil {
		return nil, err
	}

	return &User{
		UserID:   userID.String(),
		Email:    u.Email,
		Password: password,
		Payload:  u.Payload,
	}, nil
}
