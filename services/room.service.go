package services

import (
	"context"
	"errors"
	"github.com/dchest/uniuri"
	"go-pokerchips/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RoomService interface {
	CreateRoom(*models.RoomInput) (*models.DBRoom, error)
	FindRoomByUri(uri string) (*models.DBRoom, error)
}

type RoomServiceImpl struct {
	collection *mongo.Collection
}

func NewRoomService(collection *mongo.Collection) RoomService {
	return &RoomServiceImpl{collection}
}

func (rs *RoomServiceImpl) CreateRoom(room *models.RoomInput) (*models.DBRoom, error) {

	ctx := context.Background()

	room.Uri = uniuri.NewLen(5)
	room.CreatedAt = time.Now()
	room.UpdatedAt = room.CreatedAt

	res, err := rs.collection.InsertOne(ctx, room)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("room already exists")
		}
		return nil, err
	}

	index := mongo.IndexModel{Keys: bson.M{"uri": 1}, Options: options.Index().SetUnique(true)}

	if _, err := rs.collection.Indexes().CreateOne(ctx, index); err != nil {
		return nil, errors.New("could not create index for title")
	}

	var newRoom *models.DBRoom
	query := bson.M{"_id": res.InsertedID}

	err = rs.collection.FindOne(ctx, query).Decode(&newRoom)
	if err != nil {
		return nil, err
	}

	return newRoom, nil
}

func (rs *RoomServiceImpl) FindRoomByUri(uri string) (*models.DBRoom, error) {

	ctx := context.TODO()
	var room *models.DBRoom

	query := bson.M{"uri": uri}
	err := rs.collection.FindOne(ctx, query).Decode(&room)

	if err != nil {
		return nil, err
	}

	return room, nil
}
