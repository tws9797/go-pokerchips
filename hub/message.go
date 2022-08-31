package hub

import (
	"encoding/json"
	"log"
)

const (
	AddPot            = "add-pot"
	RetrievePot       = "retrieve-pot"
	SendMessageAction = "send-message"
	JoinRoomAction    = "join-room"
	LeaveRoomAction   = "leave-room"
)

type Message struct {
	Action  string  `json:"action"`
	Message string  `json:"message"`
	Pot     int     `json:"pot,omitempty"`
	Sender  *Client `json:"sender"`
}

func (message *Message) encode() []byte {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return msg
}
