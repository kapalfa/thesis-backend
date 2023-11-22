package main 
import (
	"github.com/gofiber/fiber/v2"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/routes"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os/exec"
)

func main(){
	//run the script to create localhost.pem and localhost-key.pem
	//git problems
	// fmt.Println
	cmd := exec.Command("sh", "-c", "sh ./createKey.sh")
	err := cmd.Run()
	if err != nil {	
		log.Fatalf("Failed to run script: %v\n", err)
	}

	database.ConnectDB()
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: "https://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, PATCH",
	}))
	routes.Setup(app)
	if err := app.ListenTLS(":8443", "./mkcert/localhost.pem", "./mkcert/localhost-key.pem"); err != nil {
		log.Fatalf("Error listen: %v", err)
	}
	
}