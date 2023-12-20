package filesCRUD 

import (
	"github.com/gofiber/fiber/v2"
	"os"
	"path/filepath"
	"strings"
)
type Dir struct {
	Name string `json:"name"`
	IsDir bool `json:"isDir"`
	Filepath string `json:"filepath,omitempty"`
	Children map[string]*Dir `json:"children"`
}

// get project directory 
func GetFiles(c *fiber.Ctx) error { 
	id := c.Params("id")
	dirPath := "./uploads/" + id
	root := &Dir{Name: id, IsDir: true, Children: make(map[string]*Dir), Filepath: dirPath}
	
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get files", "data": err})
	}
	if len(entries) == 0 { // if the directory is empty
		return c.JSON(root)
	}

	err = filepath.WalkDir(dirPath, func(path string, info os.DirEntry, err error) error {
		if err != nil {
			return err
		}
	
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return err
		}

		currentDir := root
		parts := strings.Split(relPath, string(os.PathSeparator))
		for i, part := range parts {
			if part == "." || part == ".." {
				continue
			}
			if _, ok := currentDir.Children[part]; !ok {
				newDir := &Dir{Name: part, IsDir: info.IsDir(), Children: make(map[string]*Dir)}
				newDir.Filepath = path
				currentDir.Children[part] = newDir
			}
			currentDir = currentDir.Children[part]

			if i == len(parts)-1 && !info.IsDir() {
				break
			}
		}
		return nil
	})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Couldn't get files", "data": err})
	}	
	return c.JSON(root)
}
