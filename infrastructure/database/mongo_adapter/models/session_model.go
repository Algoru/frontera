package mongomodels

import (
	"github.com/Algoru/frontera/domain/entity"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Credential struct {
	ID        primitive.ObjectID `bson:"_id"`
	UserID    string             `bson:"user_id"`
	Token     string             `bson:"token"`
	CreatedAt time.Time          `bson:"created_at"`
	ExpiresAt time.Time          `bson:"expires_at"`
}

func (c *Credential) ToEntity() entity.Credential {
	return entity.Credential{
		UserID:    c.UserID,
		Token:     c.Token,
		ExpiresAt: c.ExpiresAt,
	}
}
