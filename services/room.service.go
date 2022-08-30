package services

import (
	"context"
	"errors"
	"go-pokerchips/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type RoomService interface {
	CreateRoom(room *models.RoomInput) (*models.DBRoom, error)
	FindRoomByName(name string) (*models.DBRoom, error)
}

type RoomServiceImpl struct {
	collection *mongo.Collection
}

func NewRoomService(collection *mongo.Collection) RoomService {
	return &RoomServiceImpl{collection}
}

func (rs *RoomServiceImpl) CreateRoom(room *models.RoomInput) (*models.DBRoom, error) {

	ctx := context.Background()

	room.CreatedAt = time.Now()
	room.UpdatedAt = room.CreatedAt

	res, err := rs.collection.InsertOne(ctx, room)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("room name already exist")
		}
		return nil, err
	}

	var newRoom *models.DBRoom
	query := bson.M{"_id": res.InsertedID}

	err = rs.collection.FindOne(ctx, query).Decode(&newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (rs *RoomServiceImpl) FindRoomByName(name string) (*models.DBRoom, error) {

	ctx := context.TODO()
	var room *models.DBRoom

	query := bson.M{"name": name}
	err := rs.collection.FindOne(ctx, query).Decode(&room)

	if err != nil {
		return nil, err
	}

	return room, nil
}
