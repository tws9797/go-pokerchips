package hub

import (
	"fmt"
	"go-pokerchips/models"
)

type Hub struct {

	// Track the rooms available
	rooms map[*Room]bool
}

func NewHub() *Hub {

	return &Hub{
		rooms: make(map[*Room]bool),
	}
}

func (hub *Hub) FindRoomByUri(uri string) *Room {

	for room := range hub.rooms {
		if room.Uri == uri {
			return room
		}
	}

	return nil
}

// CreateRoom creates room in memory and assign its pointer to hub map.
func (hub *Hub) CreateRoom(room *models.DBRoom) *Room {

	hubRoom := NewRoom(room)
	go hubRoom.RunRoom()
	hub.rooms[hubRoom] = true

	fmt.Println(hub.rooms)
	return hubRoom
}
