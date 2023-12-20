package uploads

import (
	"github.com/gofiber/fiber/v2"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func FolderUpload(c *fiber.Ctx) error {
	projectId := c.FormValue("projectId")
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on folder upload", "data": err})
	}
	files := form.File["files"]
	for _, file := range files {
		var filename string
		contentDisposition := file.Header.Get("Content-Disposition")
		parts := strings.Split(contentDisposition, ";")
		for _, part := range parts {
			if strings.Contains(part, "filename") {
				filename = strings.Split(part, "=")[1]
				filename = strings.Trim(filename, "\"")
				break
			}
		}
	 	dir := filepath.Dir(filename)
		if _, err := os.Stat("./uploads/" + projectId + "/" + dir); os.IsNotExist(err) {
	 		os.MkdirAll("./uploads/" + projectId + "/" + dir, 0755)
	 	}
	 	fileContent, err := file.Open()
	 	if err != nil {
	 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on folder upload", "data": err})
	 	}
		defer fileContent.Close()

	 	content, err := io.ReadAll(fileContent)
	 	if err != nil {
	 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on folder upload", "data": err})
		}
	 	err = os.WriteFile("./uploads/" + projectId + "/" + dir + "/" + file.Filename, content, 0644)
	 	if err != nil {
	 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on folder upload", "data": err})
	 	}
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Folder uploaded", "data": nil})
}