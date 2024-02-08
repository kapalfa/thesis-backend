package chat 

import (
	"log"
)

type Room struct {
	id 			string
	clients 	map[*Client]bool
	broadcast 	chan *Message
	register 	chan *Client
	unregister 	chan *Client
}

func NewRoom(id string) *Room {
	return &Room{
		id: id,
		clients: make(map[*Client]bool),
		broadcast: make(chan *Message),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (room *Room) Run() {
	for {
		select {
		case client := <-room.register:
			room.registerClientInRoom(client)
			log.Println("New client registered", client)
		case client := <-room.unregister:
			room.unregisterClientInRoom(client)
			log.Println("Client unregistered", client)
		case message := <-room.broadcast:
			for client := range room.clients {
				client.send <- message.encode()
			}
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client] = true
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client]; ok {
		delete(room.clients, client)
	}
}

func (room *Room) broadcastMessage(message []byte) {
	for client := range room.clients {
		client.send <- message
	}
}