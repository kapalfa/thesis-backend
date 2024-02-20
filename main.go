package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/controllers/chat"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/routes"
)

func main() {
	database.ConnectDB()
	config.ConfigStorage()
	r := mux.NewRouter()
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	allowedOrigins := []string{os.Getenv("FRONTEND_URL")}
	origins := handlers.AllowedOrigins(allowedOrigins) // env
	log.Println("FRONTEND_URL: ", os.Getenv("FRONTEND_URL"))
	//origins := handlers.AllowedOrigins([]string{"*"})
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
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//err = http.ListenAndServeTLS(":"+port, "./mkcert/localhost.pem", "./mkcert/localhost-key.pem", r)
	err = http.ListenAndServe(":"+port, r)
	if err != nil && err != io.EOF {
		log.Fatalf("Error listen: %v", err)
	}
}
