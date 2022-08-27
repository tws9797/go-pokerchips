package services

import (
	"go-pokerchips/models"
)

type UserService interface {
	RemoveUser(string) error
	GetAllUsers() ([]*models.DBUser, error)
	FindUserById(string) (*models.DBUser, error)
	FindUserByUsername(string) (*models.DBUser, error)
}
