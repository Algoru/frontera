package mongoadapter

import (
	"context"
	mongomodels "github.com/Algoru/frontera/infrastructure/database/mongo_adapter/models"
	"time"

	"github.com/Algoru/frontera/domain/entity"
	"go.mongodb.org/mongo-driver/bson"
)

const sessionsCollectionName = "sessions"

// GetCredentialByToken
func (ma *MongoAdapter) GetCredentialByToken(token string) (*entity.Credential, error) {
	result := ma.client.Database(ma.Database).Collection(sessionsCollectionName).FindOne(context.TODO(), bson.M{"token": token})
	if err := result.Err(); err != nil {
		return nil, err
	}

	credential := mongomodels.Credential{}
	if err := result.Decode(&credential); err != nil {
		return nil, err
	}

	credentialEntity := credential.ToEntity()
	return &credentialEntity, nil
}

// AddUserSession
func (ma *MongoAdapter) AddUserSession(c *entity.Credential) error {
	_, err := ma.client.Database(ma.Database).Collection(sessionsCollectionName).InsertOne(context.TODO(), bson.M{
		"user_id":    c.UserID,
		"token":      c.Token,
		"expires_at": c.ExpiresAt,
		"created_at": time.Now(),
	})

	return err
}

// RemoveUserSessions
func (ma *MongoAdapter) RemoveUserSessions(userID string) error {
	_, err := ma.client.Database(ma.Database).Collection(sessionsCollectionName).
		DeleteMany(context.TODO(), bson.M{"user_id": userID})
	return err
}

// RemoveSingleSession
func (ma *MongoAdapter) RemoveSingleSession(userID, token string) error {
	_, err := ma.client.Database(ma.Database).Collection(sessionsCollectionName).
		DeleteMany(context.TODO(), bson.M{"user_id": userID, "token": token})
	return err
}
