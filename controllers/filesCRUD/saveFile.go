package filesCRUD 

import (
 	"github.com/gofiber/fiber/v2"
	"os"	
)

func SaveFile(c *fiber.Ctx) error {
	path := c.Params("*")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "File does not exist"})
	}	
	
	  // Parse the multipart form:
	//if form, err := c.MultipartForm(); err == nil {
	//	file := form.File["file"][0]
	//	fmt.Println("file", file)
	//}
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	}
	
	err = c.SaveFile(file, path)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	}
	//err = c.SaveFile(file, path + "/" + filename)
	//if err != nil {
	//	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	//}
	//var input map[string]interface{}
	//err := c.BodyParser(&input)
	//if err != nil {
	//	return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't save file", "data": err})
	//}
	// file, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't save file", "data": err})
	// }
	// defer file.Close()

	// err = os.WriteFile(path, input, 0644)
	//encoder := json.NewEncoder(file)
	//err = encoder.Encode(input)
	// if err != nil {
	// 	return c.Status(400).JSON(fiber.Map{"status": "error", "message": "Couldn't save file", "data": err})
	// }
	return c.JSON(fiber.Map{"status": "success", "message": "Saved file"})
}