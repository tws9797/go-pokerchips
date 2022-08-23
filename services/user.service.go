package services

import (
	"go-pokerchips/models"
)

type UserService interface {
	FindUserById(string) (*models.UserDBResponse, error)
	FindUserByUsername(string) (*models.UserDBResponse, error)
}
