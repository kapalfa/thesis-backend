package chat

type WsServer struct {
	rooms      map[string]*Room
	register   chan *Client
	unregister chan *Client
	clients    map[*Client]bool
}

func NewWebsocketServer() *WsServer {
	return &WsServer{
		rooms:      make(map[string]*Room),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (server *WsServer) Run() {
	for {
		select {

		case client := <-server.register:
			server.registerClient(client)

		case client := <-server.unregister:
			server.unregisterClient(client)
		}
	}
}

func (server *WsServer) GetRoom(chatId string) *Room {
	room, ok := server.rooms[chatId]
	if !ok {
		room = server.createRoom(chatId)
		server.rooms[chatId] = room
		go room.Run()
	}
	return room
}
func (server *WsServer) createRoom(chatId string) *Room {
	room := NewRoom(chatId)
	go room.Run()
	server.rooms[chatId] = room
	return room
}
func (server *WsServer) registerClient(client *Client) {
	server.clients[client] = true
}
func (server *WsServer) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
	}
}
func (room *Room) GetId() string {
	return room.id
}
