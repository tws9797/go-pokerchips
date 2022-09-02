package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/dchest/uniuri"
	"go-pokerchips/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type RoomService interface {
	CreateRoom(*models.CreateRoomInput) (*models.DBRoom, error)
	FindRoomByUri(string) (*models.DBRoom, error)
	RegisterUserInRoom(string, string) error
	AddPot(string, string, int) (*models.UpdatePotResponse, error)
	RetrievePot(string, string, int) (*models.UpdatePotResponse, error)
}

type RoomServiceImpl struct {
	collection *mongo.Collection
}

func NewRoomService(collection *mongo.Collection) RoomService {
	return &RoomServiceImpl{collection}
}

func (rs *RoomServiceImpl) CreateRoom(room *models.CreateRoomInput) (*models.DBRoom, error) {

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

func (rs *RoomServiceImpl) FindRoomById(id string) (*models.DBRoom, error) {

	objId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	var room *models.DBRoom

	query := bson.M{"_id": objId}
	err = rs.collection.FindOne(ctx, query).Decode(&room)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (rs *RoomServiceImpl) FindRoomByUri(uri string) (*models.DBRoom, error) {

	ctx := context.Background()
	var room *models.DBRoom

	query := bson.M{"uri": uri}
	err := rs.collection.FindOne(ctx, query).Decode(&room)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func (rs *RoomServiceImpl) RegisterUserInRoom(id string, name string) error {

	ctx := context.Background()

	room, err := rs.FindRoomById(id)
	if err != nil {
		return err
	}

	if _, ok := room.Record[name]; !ok {
		room.Record[name] = 1000
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{"_id", obId}}
	update := bson.D{{"$set", bson.D{
		{"record", room.Record},
	}}}

	_, err = rs.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return err
	}

	return nil
}

func (rs *RoomServiceImpl) AddPot(id string, name string, chips int) (*models.UpdatePotResponse, error) {

	ctx := context.Background()

	fmt.Println(id)
	fmt.Println(name)
	fmt.Println(chips)

	room, err := rs.FindRoomById(id)
	if err != nil {
		return nil, err
	}

	if _, ok := room.Record[name]; ok {
		if room.Record[name] < chips {
			return nil, err
		}
		room.Record[name] -= chips
	}
	room.Pot += chips

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{"_id", obId}}
	update := bson.D{{"$set", bson.D{
		{"pot", room.Pot},
		{"record", room.Record},
	}}}

	_, err = rs.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return nil, err
	}

	updatePotResp := &models.UpdatePotResponse{
		Pot:          room.Pot,
		CurrentChips: room.Record[name],
		Sender:       name,
	}
	
	return updatePotResp, err
}

func (rs *RoomServiceImpl) RetrievePot(id string, name string, chips int) (*models.UpdatePotResponse, error) {
	ctx := context.Background()

	room, err := rs.FindRoomById(id)
	if err != nil {
		return nil, err
	}

	room.Pot -= chips
	if room.Pot < 0 {
		return nil, errors.New("pot is not enough ")
	}
	if _, ok := room.Record[name]; ok {
		room.Record[name] += chips
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{"_id", obId}}
	update := bson.D{{"$set", bson.D{
		{"pot", room.Pot},
		{"record", room.Record},
	}}}

	_, err = rs.collection.UpdateOne(ctx, query, update)

	if err != nil {
		return nil, err
	}

	var updatePotResp *models.UpdatePotResponse
	updatePotResp.Pot = room.Pot
	updatePotResp.CurrentChips = room.Record[name]
	updatePotResp.Sender = name

	return updatePotResp, err
}
