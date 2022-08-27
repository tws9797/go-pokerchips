package services

import (
	"errors"
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

func (us *UserServiceImpl) GetAllUsers() ([]*models.UserDBResponse, error) {
	ctx := context.TODO()

	var users []*models.UserDBResponse

	cur, err := us.collection.Find(ctx, bson.D{})

	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		user := &models.UserDBResponse{}

		err = cur.Decode(user)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if len(users) == 0 {
		return []*models.UserDBResponse{}, nil
	}

	return users, nil
}

func (us *UserServiceImpl) RemoveUser(id string) error {
	ctx := context.TODO()

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := us.collection.DeleteOne(ctx, query)

	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil

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
