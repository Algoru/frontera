package mongomodels

import (
	"time"

	"github.com/Algoru/frontera/domain/entity"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	Email     string             `bson:"email"`
	Password  string             `bson:"password"`
	Payload   bson.M             `bson:"payload"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func (u *User) ToEntity() (*entity.User, error) {
	userID, err := uuid.Parse(u.UserID)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		UserID:    userID,
		Email:     u.Email,
		Password:  u.Password,
		Payload:   u.Payload,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}, nil
}
