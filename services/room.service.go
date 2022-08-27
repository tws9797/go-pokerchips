package services

import "go-pokerchips/models"

type RoomService interface {
	AddRoom(room *models.RoomInput) (*models.RoomDBResponse, error)
	FindRoomByName(name string) (*models.RoomDBResponse, error)
}
