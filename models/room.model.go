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
	Name      string             `json:"name" bson:"name"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

type RoomInput struct {
	Name      string    `json:"name" bson:"name" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

func (room *DBRoom) GetName() string {
	return room.Name
}

func (room *DBRoom) GetID() string {
	return room.ID.Hex()
}
