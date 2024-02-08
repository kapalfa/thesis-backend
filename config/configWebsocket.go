package config 

import (
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"github.com/gorilla/websocket"
	"github.com/creack/pty"
	"os"
	"sync"
	"bufio"
	"github.com/gorilla/mux"
)
type shell struct {
	cmd *exec.Cmd
	tty *os.File
}
type Message struct {
	LineContent string `json:"lineContent"`
}
type User struct {
	conn *websocket.Conn
	mu sync.Mutex
	sh *shell
}
type UsersManager struct {
	users map[string]*User
	mu sync.Mutex
	connections int
}
var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}
var mutex = &sync.Mutex{}
var um = NewUserManager()
func startShell(projID string) *shell {
	log.Println("Starting shell")
	cmd := exec.Command("bash")
	cmd.Dir = "./uploads/" + projID
	tty, err := pty.Start(cmd)
	if err != nil {
		log.Println("Error starting pty: ", err)
		return nil
	}
	return &shell{cmd, tty}
}
func NewUserManager() *UsersManager {
	return &UsersManager{
		users: make(map[string]*User),
	}
}
func (um *UsersManager) addUser(id string, projid string, conn *websocket.Conn) {
	um.mu.Lock()
	defer um.mu.Unlock()
	if user, ok := um.users[id]; ok {
		err := user.conn.Close()
		if err != nil {
			log.Println("Error closing WebSocket connection: ", err)
		}
		um.connections--
	}
	sh := startShell(projid)
	if sh == nil {
		log.Println("Error starting shell")
		return
	}
	user := &User{conn: conn, sh: sh}
	um.users[id] = user
	um.connections++
}
func (um *UsersManager) removeUser(id string) {
	um.mu.Lock()
	defer um.mu.Unlock()
	if _, ok := um.users[id]; ok {
		delete(um.users, id)
		um.connections--
	}
}
func (um *UsersManager) getUser(id string) *User {
	um.mu.Lock()
	defer um.mu.Unlock()
	if user, ok := um.users[id]; ok {
		return user
	}
	return nil
}
func (um *UsersManager) getOpenConnections() int {
	um.mu.Lock()
	defer um.mu.Unlock()
	return um.connections
}
func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request){
	// userid := path.Base(r.URL.Path)
	// if userid == "" {
	// 	log.Println("No user id provided")
	// 	return
	// }
	// id, err := strconv.ParseUint(userid, 10, 32)
	// if err != nil {
	// 	log.Println("Error parsing user id: ", err)
	// 	return
	// }
	vars := mux.Vars(r)
	projID := vars["projId"]
	userID := vars["userId"]

	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket: ", err)
		return
	}
	defer conn.Close()
	um.addUser(userID, projID, conn)
	defer um.removeUser(userID)

	user:= um.getUser(userID)
	if user == nil {
		log.Println("Error getting user")
		return
	}
	log.Println("connections: ", um.getOpenConnections())
	quit := make(chan struct{})
	defer func() {
	 	log.Println("Shell closed")
	 	err = user.sh.cmd.Process.Kill()
	 	if err != nil {
	 		log.Println("Error killing process: ", err)
	 	}
	 	_, err = user.sh.cmd.Process.Wait()
		if err != nil {
	 		log.Println("Error waiting for process: ", err)
	 	}
	 	err = user.sh.tty.Close()
	 	if err != nil {
	 		log.Println("Error closing tty: ", err)
	 	}
	 	user.sh = nil
	 	quit <- struct{}{} // send signal to quit goroutine
	}()
	go func() {
		scanner := bufio.NewScanner(user.sh.tty)
	 	for { 
			select {
			case <-quit:
				log.Println("Quit channel received")
				return
			default:
				if scanner.Scan() {
					log.Println("Read from pty/cmd: ", scanner.Text())
					err = conn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
					if err != nil {
						log.Println("Error writing to WebSocket: ", err)
						err = conn.Close()
						if err != nil {
							log.Println("Error closing WebSocket: ", err)
						}
						conn = nil
						return
					}
				}
			}
		}
	}()
	for {
		_, p, err := conn.ReadMessage()
		if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
			log.Println("Websocket closed from client: ", err)
			return			
		}
					
		var msg Message
		err = json.Unmarshal(p, &msg)
		if err != nil {
			log.Println("Error unmarshalling WebSocket message: ", err)
			return
		}
		mutex.Lock()
		_, err = user.sh.tty.WriteString(msg.LineContent+"\n")
		if err != nil {
			log.Println("Error writing to tty: ", err)
			return
		}
		mutex.Unlock()
	}
}