package chat

import (
	"encoding/json"
	"log"

	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
)

const sendMessageAction = "sendMessage"
const joinRoomAction = "joinRoom"
const leaveRoomAction = "leaveRoom"

type Message struct {
	Action   string `json:"action"`
	Message  string `json:"message"`
	RoomId   string `json:"roomId"`
	SenderId string `json:"senderId"`
}

type response struct {
	Message  string `json:"message"`
	SenderId string `json:"senderId"`
	Email    string `json:"email"`
}

func (message *Message) encode(room *Room) []byte {
	var user models.User
	var res response
	res.Message = message.Message
	res.SenderId = message.SenderId

	if message.Action == sendMessageAction {
		database.DB.Where("id = ?", message.SenderId).First(&user)
		res.Email = user.Email
	}
	encoded, err := json.Marshal(res)
	if err != nil {
		log.Println("Error encoding message", err)
	}
	return encoded
}
