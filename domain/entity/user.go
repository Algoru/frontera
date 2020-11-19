package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
