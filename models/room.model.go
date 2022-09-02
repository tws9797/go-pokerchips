package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type DBRoom struct {
	Id        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Uri       string             `json:"uri" bson:"uri"`
	Pot       int                `json:"pot" bson:"pot"`
	Record    map[string]int     `json:"record" bson:"record"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
}

type CreateRoomInput struct {
	Creator   string         `json:"name" bson:"name"`
	Uri       string         `json:"uri" bson:"uri"`
	Record    map[string]int `json:"record" bson:"record"`
	CreatedAt time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt" bson:"updatedAt"`
}

type JoinRoomInput struct {
	User string `json:"name"`
	Uri  string `json:"uri"`
}

type UpdatePotResponse struct {
	Pot          int    `json:"pot"`
	Sender       string `json:"name"`
	CurrentChips int    `json:"currentChips"`
}
