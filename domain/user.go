package domain

import (
	"strings"

	"github.com/microcosm-cc/bluemonday"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User struct represents an user inside a MongoDB
// database. It contains required fields such an email and
// a password. The Payload attribute contains what ever
// platform needs to store about their users.
type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   string             `bson:"user_id"`
	Email    string             `bson:"email"`
	Password string             `bson:"password"`
	Payload  bson.M             `bson:"payload"`
}

// Sanitize
func (u *User) Sanitize() {
	u.Email = strings.ToLower(u.Email)

	pol := bluemonday.StrictPolicy()
	u.Email = pol.Sanitize(u.Email)
}
