package filesCRUD

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"os"
)

func GetFile(c* fiber.Ctx) error {
	path := c.Params("*")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "File does not exist"})
	}
	fmt.Println("path", path)
	//path := c.Query("path")

	err := c.SendFile(path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Could not read file", "data": err})
	}

	return nil
}