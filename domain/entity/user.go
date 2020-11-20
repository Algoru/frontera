package entity

import (
	"strings"
	"time"

	"github.com/Algoru/frontera/configuration"
	"github.com/google/uuid"
	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

// User struct represents an user inside a MongoDB
// database. It contains required fields such an email and
// a password. The Payload attribute contains what ever
// platform needs to store about their users.
type User struct {
	ID        primitive.ObjectID `json:"-"`
	UserID    uuid.UUID          `json:"user_id"`
	Email     string             `json:"email"`
	Password  string             `json:"password,omitempty"`
	Payload   bson.M             `json:"payload"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

// RemoveSensible removes sensible data from User struct such as passwod
func (u *User) RemoveSensible() {
	u.Password = ""
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
		UserID:   userID,
		Email:    u.Email,
		Password: password,
		Payload:  u.Payload,
	}, nil
}
