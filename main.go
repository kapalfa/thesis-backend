package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/controllers/chat"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/routes"
)

func main() {
	database.ConnectDB()
	config.ConfigStorage()
	r := mux.NewRouter()

	allowedOrigins := []string{os.Getenv("FRONTEND_URL")}
	origins := handlers.AllowedOrigins(allowedOrigins)
	log.Println("FRONTEND_URL: ", os.Getenv("FRONTEND_URL"))
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"})
	headers := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization"})
	credentials := handlers.AllowCredentials()

	r.Use(handlers.CORS(origins, methods, headers, credentials))
	routes.Setup(r)

	wsServer := chat.NewWebsocketServer()
	go wsServer.Run()

	r.HandleFunc("/sockets/chat/{id}", func(w http.ResponseWriter, r *http.Request) {
		wsServer.CreateConnections(w, r)
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := http.ListenAndServe(":"+port, r)
	if err != nil && err != io.EOF {
		log.Fatalf("Error listen: %v", err)
	}
}
