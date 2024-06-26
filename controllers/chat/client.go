package chat

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	rooms    map[*Room]bool
	wsServer *WsServer
	conn     *websocket.Conn
	send     chan []byte
}

func NewClient(conn *websocket.Conn, r *Room, wsServer *WsServer) *Client {
	return &Client{
		rooms:    make(map[*Room]bool),
		wsServer: wsServer,
		conn:     conn,
		send:     make(chan []byte, 256),
	}
}

func (client *Client) handleJoinRoomMessage(message Message) *Room {
	roomId := message.RoomId
	room := client.wsServer.GetRoom(roomId)
	if room == nil {
		room = client.wsServer.createRoom(roomId)
	}
	if !client.isInRoom(room) {
		client.rooms[room] = true
		room.register <- client
	}
	return room
}
func (client *Client) handleLeaveRoomMessage(message Message) {
	roomId := message.RoomId
	room := client.wsServer.GetRoom(roomId)
	if room == nil {
		return
	}
	room.unregister <- client
}
func (client *Client) handleNewMessage(jsonMessage []byte) {
	var message Message

	err := json.Unmarshal(jsonMessage, &message)
	if err != nil {
		log.Println("Error unmarshalling message", err)
		return
	}

	switch message.Action {
	case sendMessageAction:
		roomId := message.RoomId
		if room := client.wsServer.GetRoom(roomId); room != nil {
			room.broadcast <- &message
		}
	case joinRoomAction:
		client.handleJoinRoomMessage(message)
	case leaveRoomAction:
		client.handleLeaveRoomMessage(message)
	}
}

func (client *Client) isInRoom(room *Room) bool {
	if _, ok := client.rooms[room]; ok {
		return true
	}
	return false
}
