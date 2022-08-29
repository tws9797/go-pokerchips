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
			room.unregisterClient(client)
		case message := <-room.broadcast:
			fmt.Println("broadcast~~~")
			fmt.Println(message)
			room.broadcastClientsInRoom(message.encode())
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {

	fmt.Println("registerClientInRoom")
	// Sending message before register the user so the user won't see the message him/herself joined
	room.notifyClientJoined(client)
	room.clients[client] = true
}

func (room *Room) unregisterClient(client *Client) {

	fmt.Println("unregisterClient")
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastClientsInRoom(message []byte) {

	fmt.Println("broadcastClientsInRoom")
	fmt.Println(room.clients)
	fmt.Println(string(message))
	for client := range room.clients {
		client.send <- message
	}
}

func (room *Room) notifyClientJoined(client *Client) {

	fmt.Println("notifyClientJoined")
	message := &Message{
		Action:  SendMessageAction,
		Target:  room,
		Message: fmt.Sprintf(welcomeMessage, client.name),
	}
	fmt.Println(message)

	room.broadcastClientsInRoom(message.encode())
}
