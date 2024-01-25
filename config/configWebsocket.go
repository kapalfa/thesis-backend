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
)
type shell struct {
	cmd *exec.Cmd
	tty *os.File
}
type Message struct {
	LineContent string `json:"lineContent"`
	Projectid string	`json:"id"`
}
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}
var messageReceived = make(chan int)
var globalConn *websocket.Conn
var currentShell *shell
var mutex = &sync.Mutex{}
func startShell() *shell {
	log.Println("Starting shell")
	cmd := exec.Command("bash")
	cmd.Dir = "./uploads"
	tty, err := pty.Start(cmd)
	if err != nil {
		log.Println("Error starting pty: ", err)
		return nil
	}
	return &shell{cmd, tty}
}
func initWebSocketConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading connection to WebSocket: ", err)
		return nil, err 
	}
	globalConn = conn
	log.Println("Websocket connection established")
	return conn, nil
}
func HandleWebsocketConnection(w http.ResponseWriter, r *http.Request){
	var err error
	quit := make(chan struct{})
	log.Println("Handling WebSocket connection", globalConn)

	if globalConn == nil {
		log.Println("Creating new WebSocket connection")
	 	globalConn, err = initWebSocketConnection(w, r)
	 	if err != nil {
			log.Println("Cant't create websocket connection")
	 		return
		}
	/*	defer func() {
	//		log.Println("Closing WebSocket connection")
	//		err := globalConn.Close()
	//		if err != nil {
	//			log.Println("Error closing WebSocket connection: ", err)
	//		}	
	//		log.Println("WebSocket connection closed")
	//		globalConn = nil
		}()*/
	}
	if currentShell == nil {

		currentShell = startShell()
		defer func() {
			log.Println("Shell closed")
			err = currentShell.cmd.Process.Kill()
			if err != nil {
				log.Println("Error killing process: ", err)
			}
			_, err = currentShell.cmd.Process.Wait()
			if err != nil {
				log.Println("Error waiting for process: ", err)
			}
			err = currentShell.tty.Close()
			if err != nil {
				log.Println("Error closing tty: ", err)
			}
			currentShell = nil
			quit <- struct{}{} // send signal to quit goroutine
		}()
	}
	go func() {
		scanner := bufio.NewScanner(currentShell.tty)
	 	for { 
			select {
			case <-quit:
				log.Println("Quit channel received")
				return
			default:
	 //	case <-messageReceived:
				if scanner.Scan() {
					log.Println("Read from pty/cmd: ", scanner.Text())
					err = globalConn.WriteMessage(websocket.TextMessage, []byte(scanner.Text()))
					if err != nil {
						log.Println("Error writing to WebSocket: ", err)
						err = globalConn.Close()
						if err != nil {
							log.Println("Error closing WebSocket: ", err)
						}
						globalConn = nil
						log.Println("WebSocket closed: ", globalConn)
						return
					}
				}
			}
		}
	//	}
	 		//	n, err := currentShell.tty.Read(buf) // read from pty/cmd
		 	//	if n != 0 { // don't write empty messages
		 	//		err = globalConn.WriteMessage(websocket.TextMessage, buf[:n])
	}()
	for {
		_, p, err := globalConn.ReadMessage()
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
		_, err = currentShell.tty.WriteString(msg.LineContent+"\n")
		if err != nil {
			log.Println("Error writing to tty: ", err)
			return
		}
		mutex.Unlock()
	
	//	messageReceived <- true

		//cmd := exec.Command("bash", "-c", "cd "+dir+" && "+msg.LineContent)
	//	cmd := exec.Command("bash", "-c", msg.LineContent)
	//	tty, err := pty.Start(cmd)
		// buf := make([]byte, 1024)
		// n, err := currentShell.tty.Read(buf)
		// err = globalConn.WriteMessage(websocket.TextMessage, buf[:n])
	// 		err = currentShell.cmd.Process.Signal(syscall.SIGTERM)
	}
}