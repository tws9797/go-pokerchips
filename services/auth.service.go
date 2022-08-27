package services

import (
	"go-pokerchips/models"
)

type AuthService interface {
	SignUpUser(*models.SignUpInput) (*models.DBUser, error)
}
