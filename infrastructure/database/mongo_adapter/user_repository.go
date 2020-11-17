package mongoadapter

import (
	"context"

	"github.com/Algoru/frontera/domain/entity"
	userrepository "github.com/Algoru/frontera/repository/user_repository"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const usersCollectionName = "users"

// CreateUser
func (ma *MongoAdapter) CreateUser(u *userrepository.User) (*entity.User, error) {
	result, err := ma.client.Database(ma.Database).Collection(usersCollectionName).InsertOne(context.TODO(), bson.M{
		"user_id":    u.UserID,
		"email":      u.Email,
		"password":   u.Password,
		"payload":    u.Payload,
		"created_at": u.CreatedAt,
		"updated_at": u.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	u.ID = result.InsertedID.(primitive.ObjectID)
	entityUser := u.ToRawEntity()

	return &entityUser, nil
}

// GetUser
func (ma *MongoAdapter) GetUser(userID uuid.UUID) (*entity.User, error) {
	result := ma.client.Database(ma.Database).Collection(usersCollectionName).
		FindOne(context.TODO(), bson.M{"user_id": userID.String()})
	if err := result.Err(); err != nil {
		return nil, err
	}

	userFound := userrepository.User{}
	if err := result.Decode(&userFound); err != nil {
		return nil, err
	}

	rawEntity := userFound.ToRawEntity()
	return &rawEntity, nil
}

// UpdateUser
func (ma *MongoAdapter) UpdateUser(userID uuid.UUID, u userrepository.User) (*entity.User, error) {
	// TODO (@Algoru): Identify how user update should me made
	return nil, nil
}

// DeleteUser
func (ma *MongoAdapter) DeleteUser(userID uuid.UUID) (*entity.User, error) {
	opts := &options.FindOneAndDeleteOptions{Projection: options.Before}

	result := ma.client.Database(ma.Database).Collection(usersCollectionName).
		FindOneAndDelete(context.TODO(), bson.M{"user_id": userID}, opts)
	if err := result.Err(); err != nil {
		return nil, err
	}

	userFound := entity.User{}
	if err := result.Decode(&userFound); err != nil {
		return nil, err
	}

	return &userFound, nil
}
