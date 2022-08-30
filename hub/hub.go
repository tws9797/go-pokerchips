package hub

import (
	"fmt"
	"go-pokerchips/services"
)

type Hub struct {

	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Track the rooms available
	rooms map[*Room]bool

	roomService services.RoomService
}

func NewHub() *Hub {

	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		rooms:      make(map[*Room]bool),
	}
}

// Run the websocket server to listening on register, unregister and broadcast channel
func (hub *Hub) Run() {

	for {
		select {
		case client := <-hub.register:
			hub.registerClient(client)
		case client := <-hub.unregister:
			hub.unregisterClient(client)
		case message := <-hub.broadcast:
			hub.broadcastClients(message)
		}
	}
}

// registerClient to register client pointer in the hub
func (hub *Hub) registerClient(client *Client) {

	fmt.Println("registerClient")
	hub.clients[client] = true
}

// unregisterClient remove client from the hub and close the channel
func (hub *Hub) unregisterClient(client *Client) {

	fmt.Println("unregisterClient")
	if _, ok := hub.clients[client]; ok {
		delete(hub.clients, client)
	}
}

// broadcastClients send message to all the clients in the hub
func (hub *Hub) broadcastClients(message []byte) {

	fmt.Println("broadcastClients")
	for client := range hub.clients {
		client.send <- message
	}
}

func (hub *Hub) findRoomByName(name string) *Room {

	for room := range hub.rooms {
		fmt.Println(room.name)
		if room.name == name {
			return room
		}
	}

	return nil
}

func (hub *Hub) FindRoomByID(id string) *Room {

	var foundRoom *Room
	for room := range hub.rooms {
		if room.id == id {
			return foundRoom
		}
	}

	return nil
}

func (hub *Hub) CreateRoom(name string) *Room {

	room := NewRoom(name)
	go room.RunRoom()
	hub.rooms[room] = true

	return room
}
