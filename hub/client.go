package hub

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
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

func newClient(conn *websocket.Conn, hub *Hub, name string) *Client {

	return &Client{
		name: name,
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 256),
	}
}

// writePump handles sending the messages from the hub to the websocket connection.
// A goroutine running writePump is started for each connection.
// The application ensures that there is at most one writer to a connection by executing all writes from this goroutine.
func (client *Client) writePump() {

	fmt.Println("writePump...")
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
			fmt.Println("messge writing: ")
			fmt.Println(string(message))

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

	fmt.Println("readPump...")
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

		fmt.Println("message")
		fmt.Println(string(message))

		client.handleNewMessage(message)
	}
}

func ServeWS(hub *Hub, c *gin.Context) {
	name, ok := c.Request.URL.Query()["name"]

	if !ok || len(name[0]) < 1 {
		log.Println("Url Param 'name' is missing")
		return
	}

	// Upgrade the HTTP server connection to the websocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := newClient(conn, hub, name[0])
	hub.register <- client

	go client.writePump()
	go client.readPump()
}

func (client *Client) disconnect() {
	fmt.Println("disconnect")
	fmt.Println(client)
	client.hub.unregister <- client
	client.room.unregister <- client
	client.conn.Close()
}

func (client *Client) handleNewMessage(message []byte) {
	fmt.Println(client)
	fmt.Println("handleNewMessage")
	fmt.Println(string(message))

	var msg Message
	if err := json.Unmarshal(message, &msg); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
	}

	msg.Sender = client

	switch msg.Action {
	case SendMessageAction:
		fmt.Println("sendmessageAction")
		client.room.broadcast <- &msg
	case JoinRoomAction:
		fmt.Println("JoinRoomAction")
		roomName := msg.Message
		fmt.Println(roomName)

		room := client.hub.findRoomByName(roomName)
		fmt.Println("is Room found?")
		fmt.Println(room)
		if room == nil {
			room = client.hub.createRoom(roomName)
		}
		client.room = room
		client.room.register <- client
	case LeaveRoomAction:
		fmt.Println("LeaveRoomAction")
		client.room.unregister <- client
	}
}
