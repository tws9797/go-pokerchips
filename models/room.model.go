package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Room interface {
	GetID()
	GetName()
}

type DBRoom struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Uri       string             `json:"uri" bson:"uri"`
	Pot       int                `json:"pot" bson:"pot"`
	Record    map[string]int     `json:"record" bson:"record"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type RoomInput struct {
	Creator   string         `json:"name" bson:"name"`
	Uri       string         `json:"uri" bson:"uri"`
	Record    map[string]int `json:"record" bson:"record"`
	CreatedAt time.Time      `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time      `json:"updated_at" bson:"updated_at"`
}

func (room *DBRoom) GetName() string {
	return room.Uri
}

func (room *DBRoom) GetID() string {
	return room.ID.Hex()
}
