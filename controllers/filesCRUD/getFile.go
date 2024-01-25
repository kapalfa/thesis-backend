package filesCRUD

import (
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

// func GetFile(c* fiber.Ctx) error {
// 	path := c.Params("*")
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "File does not exist"})
// 	}
// 	fmt.Println("path", path)
// 	//path := c.Query("path")

// 	err := c.SendFile(path)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Could not read file", "data": err})
// 	}

// 	return nil
// }

func GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File does not exist", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, path)
}