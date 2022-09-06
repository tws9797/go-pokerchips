package hub

import (
	"encoding/json"
	"log"
)

const (
	AddPot            = "add-pot"
	RetrievePot       = "retrieve-pot"
	UpdatePot         = "update-pot"
	JoinRoomAction    = "join-room"
	SendMessageAction = "send-message"
	LeaveRoomAction   = "leave-room"
)

type Message struct {
	Action       string `json:"action"`
	Message      string `json:"message"`
	Pot          int    `json:"pot,omitempty"`
	CurrentChips int    `json:"currentChips"`
	Sender       string `json:"sender,omitempty"`
}

func (message *Message) encode() []byte {
	msg, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return msg
}
