package main 
import (
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/routes"
	"log"
	"net/http"
	"io"
	"os/exec"
)

func main(){
	cmd := exec.Command("sh", "-c", "sh ./createKey.sh")
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to run script: %v\n", err)
	}

	database.ConnectDB()
	r := mux.NewRouter()

	origins := handlers.AllowedOrigins([]string{"https://localhost:5173"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "PATCH"})
	headers := handlers.AllowedHeaders([]string{"Origin", "Content-Type", "Accept", "Authorization"})
	credentials := handlers.AllowCredentials()

	r.Use(handlers.CORS(origins, methods, headers, credentials))

	routes.Setup(r)
	
	err = http.ListenAndServeTLS(":8443", "./mkcert/localhost.pem", "./mkcert/localhost-key.pem", r)
	if err != nil && err != io.EOF {
		log.Fatalf("Error listen: %v", err)
	}
}