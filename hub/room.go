package hub

import "fmt"

const welcomeMessage = "%s joined the room"

type Room struct {
	id   string
	name string

	//Registered clients
	clients map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client

	// Inbound messages from the clients.
	broadcast chan *Message
}

func NewRoom(name string) *Room {

	return &Room{
		name:       name,
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
			fmt.Printf("Broadcast message: %v to the room %v", message, room.name)
			room.broadcastClientsInRoom(message.encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {

	fmt.Printf("registerClientInRoom: %v\n", client.name)
	room.clients[client] = true
	room.notifyClientJoined(client)
}

func (room *Room) unregisterClientInRoom(client *Client) {

	fmt.Printf("unregisterClientInRoom: %v\n", client.name)
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
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

	message := &Message{
		Action:  SendMessageAction,
		Target:  room.name,
		Message: fmt.Sprintf(welcomeMessage, client.name),
	}
	room.broadcastClientsInRoom(message.encode())
}
