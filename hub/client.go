package hub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second

	// Max time till next pong from peer
	pongWait = 60 * time.Second

	// Send ping interval, must be less than pong wait time
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

// Client hold the connection between the websocket and hub instance.
type Client struct {
	name string

	hub *Hub

	// The websocket connection
	conn *websocket.Conn

	// Buffered channel of outbound message (message routed from a client to the end user)
	send chan []byte

	room *Room
}

func newClient(conn *websocket.Conn, hub *Hub, room *Room, name string) *Client {

	return &Client{
		conn: conn,
		hub:  hub,
		room: room,
		name: name,
		send: make(chan []byte, 256),
	}
}

// writePump handles sending the messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection.
// The application ensures that there is at most one writer to a connection by executing all writes from this goroutine.
func (client *Client) writePump() {

	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			fmt.Printf("Write Pump: %v \n", string(message))
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// readPump handles reading the messages from the websocket connection to the hub.
// The application runs readPump in a per-connection goroutine.
// The application ensures that there is at most one reader on a connection by executing all reads from this goroutine.
func (client *Client) readPump() {

	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		fmt.Printf("Read Pump Messages: %v\n", string(message))
		client.handleNewMessage(message)
	}
}

func ServeWS(hub *Hub, room *Room, name string, c *gin.Context) {

	// Upgrade the HTTP server connection to the websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, hub, room, name)

	go client.writePump()
	go client.readPump()

	room.register <- client
}

func (client *Client) disconnect() {

	fmt.Printf("%v disconnected from the room \n", client.name)
	client.room.unregister <- client
	client.conn.Close()
}

func (client *Client) handleNewMessage(message []byte) {

	fmt.Printf("handleNewMessage from %v: %v \n", client, string(message))
	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	msg.Sender = client

	switch msg.Action {
	case SendMessageAction:
		fmt.Println("SendMessageAction")
		client.room.broadcast <- &msg
	//case JoinRoomAction:
	//	fmt.Println("JoinRoomAction")
	//	client.room.register <- client
	case AddPot:
		client.addPot(msg)
	case RetrievePot:
		client.retrievePot(msg)
	case LeaveRoomAction:
		fmt.Println("LeaveRoomAction")
		client.room.unregister <- client
	}
}

func (client *Client) addPot(message Message) {

	fmt.Printf("%v trying to add pot \n", client.name)
	pot := message.Pot
	fmt.Println("pot")
	fmt.Println(pot)

	newPot, err := client.hub.roomService.AddPot(client.room.Id, client.name, pot)
	message.Message = fmt.Sprintf("%v bet for %v", client.name, pot)

	if err != nil {
		message.Message = "not enough pot to bet dude"
	}

	message.Action = "update-pot"
	message.Pot = newPot

	client.room.broadcast <- &message
}

func (client *Client) retrievePot(message Message) {

	fmt.Printf("%v trying to retrive pot \n", client.name)
	pot := message.Pot
	newPot, err := client.hub.roomService.RetrievePot(client.room.Id, client.name, pot)
	message.Message = fmt.Sprintf("%v retrieve for %v", client.name, pot)

	if err != nil {
		message.Message = "not enough pot to retrieved dude"
	}

	message.Action = "update-pot"
	message.Pot = newPot

	client.room.broadcast <- &message
}
