package chat

import(
	"net/http"
	"log"
	"github.com/gorilla/websocket"
	"github.com/gorilla/mux"
)

func (client *Client) read() {
//	defer func() {
		///client.wsServer.unregister <- client
	//	client.conn.Close()
//	}()
	defer client.conn.Close()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error: %v", err)
			return 
		}
		//client.room.broadcast <- message
		client.handleNewMessage(message)
	}
}
func (client *Client) write() {
	defer client.conn.Close()
	for msg := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
	}
}
var upgrader = &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool { return true }}
func (wsServer *WsServer) CreateConnections(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomId := vars["projId"]
	//userId := vars["id"]
	room := wsServer.GetRoom(roomId)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection: ", err)
		return
	}
	defer conn.Close()

	client := NewClient(conn, room, wsServer)
	log.Println("Client created: ", client)
	room.register <- client

	//wsServer.register <- client
	//go wsServer.Run()

	defer func() {
		room.unregister <- client
	}()

	go client.write()
	client.read()

}