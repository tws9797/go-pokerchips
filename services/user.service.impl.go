package services

import (
	"go-pokerchips/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"strings"
)

type UserServiceImpl struct {
	collection *mongo.Collection
}

func NewUserService(collection *mongo.Collection) UserService {
	return &UserServiceImpl{collection}
}

func (us *UserServiceImpl) FindUserById(id string) (*models.UserDBResponse, error) {

	ctx := context.TODO()
	oid, _ := primitive.ObjectIDFromHex(id)

	var user *models.UserDBResponse

	query := bson.M{"_id": oid}
	err := us.collection.FindOne(ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserDBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}

func (us *UserServiceImpl) FindUserByUsername(username string) (*models.UserDBResponse, error) {

	ctx := context.TODO()
	var user *models.UserDBResponse

	query := bson.M{"username": strings.ToLower(username)}
	err := us.collection.FindOne(ctx, query).Decode(&user)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.UserDBResponse{}, err
		}
		return nil, err
	}

	return user, nil
}
