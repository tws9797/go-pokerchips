package chat

// Hub maintains the set of active clients and broadcasts messages to the client
type Hub struct {
	// Keep tracks of registered clients
	clients map[*Client]bool

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client

	// Inbound messages from the clients.
	broadcast chan []byte

	// Keep tracks of the rooms in the hub
	rooms map[*Room]bool
}

// NewHub creates a new Hub type
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		rooms:      make(map[*Room]bool),
	}
}

// Run Hub server using broadcast, register and unregister channels to listen for different inbound messages
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)
		case client := <-h.unregister:
			h.unregisterClient(client)
		case message := <-h.broadcast:
			h.broadcastToClients(message)
		}
	}
}

// Register client pointer
func (h *Hub) registerClient(client *Client) {
	h.notifyClientJoined(client)
	h.listOnlineClients(client)
	h.clients[client] = true
}

// Deleting client pointer from the clients map
// Close the client's send channel to signal the client that no more messages will be sent to the client
func (h *Hub) unregisterClient(client *Client) {
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		h.notifyClientLeft(client)
	}
}

// Send read messages to registered clients
func (h *Hub) broadcastToClients(message []byte) {
	//for client := range h.clients {
	//	select {
	//	case client.send <- message:
	//	default:
	//		close(client.send)
	//		delete(h.clients, client)
	//	}
	//}
	for client := range h.clients {
		client.send <- message
	}
}

// Search the room in the hub by the name entered
func (h *Hub) findRoomByName(name string) *Room {
	var foundRoom *Room
	for room := range h.rooms {
		if room.GetName() == name {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

func (h *Hub) findRoomByID(ID string) *Room {
	var foundRoom *Room
	for room := range h.rooms {
		if room.GetId() == ID {
			foundRoom = room
			break
		}
	}

	return foundRoom
}

// Create Room
func (h *Hub) createRoom(name string, private bool) *Room {
	room := NewRoom(name, private)
	go room.RunRoom()
	h.rooms[room] = true

	return room
}

func (h *Hub) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: client,
	}

	h.broadcastToClients(message.encode())
}

func (h *Hub) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: client,
	}

	h.broadcastToClients(message.encode())
}

func (h *Hub) listOnlineClients(client *Client) {
	for existingClient := range h.clients {
		message := &Message{
			Action: UserJoinedAction,
			Sender: existingClient,
		}
		client.send <- message.encode()
	}
}

func (h *Hub) findClientByID(ID string) *Client {
	var foundClient *Client
	for client := range h.clients {
		if client.ID.String() == ID {
			foundClient = client
			break
		}
	}

	return foundClient
}
