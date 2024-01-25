package filesCRUD 

import (
	"fmt"
	"io"
	"os"
	"net/http"
	"encoding/json"
)

// func SaveFile(c *fiber.Ctx) error {
// 	path := c.Params("*")
// 	if _, err := os.Stat(path); os.IsNotExist(err) {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "File does not exist"})
// 	}	
	
// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
// 	}
// 	err = c.SaveFile(file, path)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
// 	}
	
// 	fileContent, err := file.Open()
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
// 	}
// 	defer fileContent.Close()
// 	content, err := io.ReadAll(fileContent)
// 	if err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Error on file upload", "data": err})
// 	}

// 	fmt.Println(string(content))	
// 	return c.JSON(fiber.Map{"status": "success", "message": "Saved file", "data": string(content)})
// }

func SaveFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/api/saveFile/"):]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

//	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File: ", file)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	tempFile, err := os.Create(path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer tempFile.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tempFile.Write(fileBytes)

	content, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"message": "Saved file",
		"data": string(content),
	}

	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}