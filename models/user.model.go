package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User interface {
	GetId() primitive.ObjectID
	GetName() string
}

type SignInInput struct {
	Username string `json:"username" bson:"username" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}

type SignUpInput struct {
	Username  string    `json:"username" bson:"username" binding:"required"`
	Password  string    `json:"password" bson:"password" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type CreateUserRequest struct {
	Username  string    `json:"username" bson:"username" binding:"required"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type UserDBResponse struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Username  string             `json:"username" bson:"username"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}
