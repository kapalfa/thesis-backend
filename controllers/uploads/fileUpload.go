package uploads

import (
	"github.com/gofiber/fiber/v2"
	"os"
)

//store the file inside the uploads/projectId folder
func FileUpload(c *fiber.Ctx) error {
	path := c.Params("*")
	
	if _, err := os.Stat(path); os.IsNotExist(err) { //kanonika den prepei na yparxei to path idi
		os.MkdirAll(path, 0755)
	}

	file, err := c.FormFile("file")
	if err != nil {
	 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	}
	filename := file.Filename
	err = c.SaveFile(file, path + "/" + filename)
	
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	}
	// fileContent, err := file.Open()
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	// }
	// defer fileContent.Close()
	//  content, err := io.ReadAll(fileContent)
	//  if err != nil {
	//  	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	//  }
	
	// err = os.WriteFile("./uploads/" + projectId + "/" + file.Filename, content, 0644)
	// if err != nil {
	// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
	// }
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "File uploaded"})
}