package hub

import (
	"fmt"
	"go-pokerchips/models"
	"go-pokerchips/services"
)

type Hub struct {

	// Track the rooms available
	rooms map[*Room]bool

	roomService services.RoomService
}

func NewHub(roomService services.RoomService) *Hub {

	return &Hub{
		rooms:       make(map[*Room]bool),
		roomService: roomService,
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

	hubRoom := NewRoom(hub, room)
	go hubRoom.RunRoom()
	hub.rooms[hubRoom] = true

	fmt.Println(hub.rooms)
	return hubRoom
}

func (hub *Hub) DeleteRoom(room *Room) {

	if _, ok := hub.rooms[room]; ok {
		delete(hub.rooms, room)
	}
}
