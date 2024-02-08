package main 
import (
	"io"
	"log"
	"net/http"
	"os/exec"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/routes"
	"github.com/kapalfa/go/controllers/chat"
	"os"
	"github.com/joho/godotenv"
)

func main(){
	cmd := exec.Command("sh", "-c", "sh ./createKey.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run script: %v\n", err)
	}
	
	database.ConnectDB()
	r := mux.NewRouter()
	origins := handlers.AllowedOrigins([]string{"https://localhost:5173"}) // env 
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"})
	headers := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization"})
	credentials := handlers.AllowCredentials()

	r.Use(handlers.CORS(origins, methods, headers, credentials))

	routes.Setup(r)
	
	wsServer := chat.NewWebsocketServer()
	//go room.Run()
	go wsServer.Run()

	r.HandleFunc("/sockets/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		wsServer.CreateConnections(w, r)
	})
	err = godotenv.Load()
	port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

	log.Println("Server is running on port: " + port)
	err = http.ListenAndServeTLS(":"+port, "./mkcert/localhost.pem", "./mkcert/localhost-key.pem", r)
	if err != nil && err != io.EOF {
		log.Fatalf("Error listen: %v", err)
	}
}