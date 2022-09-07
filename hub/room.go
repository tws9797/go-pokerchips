package hub

import (
	"fmt"
	"go-pokerchips/models"
)

const welcomeMessage = "> %s joined the room."
const leaveMessage = "> %s left the room."

type Room struct {
	Id     string         `json:"id"`
	Uri    string         `json:"uri"`
	Pot    int            `json:"pot"`
	Record map[string]int `json:"record"`

	hub *Hub

	//Registered clients
	clients map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client

	// Inbound messages from the clients.
	broadcast chan *Message
}

func NewRoom(hub *Hub, room *models.DBRoom) *Room {

	return &Room{
		Id:         room.Id.Hex(),
		Uri:        room.Uri,
		Pot:        room.Pot,
		Record:     room.Record,
		hub:        hub,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (room *Room) RunRoom() {

	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)
		case client := <-room.unregister:
			room.unregisterClientInRoom(client)
		case message := <-room.broadcast:
			fmt.Printf("Broadcast message: %v to the room %v", message, room.Uri)
			room.broadcastClientsInRoom(message.encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {

	fmt.Printf("registerClientInRoom: %v\n", client.name)
	room.clients[client] = true

	//Notify client with his/her username
	message := &Message{
		Pot:    room.Pot,
		Action: JoinRoomAction,
		Sender: client.name,
	}
	client.send <- message.encode()

	room.notifyClientJoined(client)
}

func (room *Room) unregisterClientInRoom(client *Client) {

	fmt.Printf("unregisterClientInRoom: %v\n", client.name)
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
		if len(room.clients) == 0 {
			room.hub.DeleteRoom(room)
			room.notifyClientJoined(client)
		}
	}
}

func (room *Room) broadcastClientsInRoom(message []byte) {

	fmt.Printf("broadcastClientsInRoom: %v\n", string(message))
	fmt.Println("The clients in room: ")
	for client := range room.clients {
		fmt.Println(client.name)
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {

	fmt.Printf("notifyClientJoined: %v\n", client.name)
	fmt.Printf("Latest Room Pot: %v\n", room.Pot)

	message := &Message{
		Action:  SendMessageAction,
		Message: fmt.Sprintf(welcomeMessage, client.name),
		Sender:  client.name,
	}

	room.broadcastClientsInRoom(message.encode())
}

func (room *Room) notifyClientLeft(client *Client) {

	fmt.Printf("notifyClientLeft: %v\n", client.name)

	message := &Message{
		Action:  SendMessageAction,
		Message: fmt.Sprintf(leaveMessage, client.name),
		Sender:  client.name,
	}
	room.broadcastClientsInRoom(message.encode())
}
