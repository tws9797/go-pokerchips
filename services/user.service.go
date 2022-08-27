package services

import (
	"go-pokerchips/models"
)

type UserService interface {
	RemoveUser(string) error
	GetAllUsers() ([]*models.UserDBResponse, error)
	FindUserById(string) (*models.UserDBResponse, error)
	FindUserByUsername(string) (*models.UserDBResponse, error)
}
