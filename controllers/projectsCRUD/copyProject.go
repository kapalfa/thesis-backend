package projectsCRUD

import (
	"net/http"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"os"
	"strconv"
	"encoding/json"
	"io"
	"path/filepath"
	"io/fs"
)
type CopyRequest struct {
	ProjectId string `json:"projectid"`
	UserId string `json:"userid"`
}
func copyDir(src string, dst string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relPath)
		if d.IsDir() {
			return os.MkdirAll(dstPath, 0777)
		}
		srcFile,err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()
		dstFile, err := os.Create(dstPath)
		if err != nil {
			return err 
		}
		defer dstFile.Close()
		_, err = io.Copy(dstFile, srcFile)
		if err != nil {
			return err
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		return os.Chmod(dstPath, info.Mode())
	})
} 
func CopyProject(w http.ResponseWriter, r *http.Request) {
	var req CopyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := &models.Project{}
	if err := database.DB.Where("id = ?", req.ProjectId).First(project).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	access := &models.Access{}
	if err := database.DB.Where("user_id = ? AND project_id = ?",req.UserId, req.ProjectId).First(access).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	newProject := &models.Project{	
		Name: project.Name,
		Description: project.Description,
		Public: false,
	}

	if err := database.DB.Create(newProject).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := os.MkdirAll("./uploads/" + strconv.Itoa(int(newProject.Id)), 0755); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := copyDir("./uploads/" + strconv.Itoa(int(project.Id)), "./uploads/" + strconv.Itoa(int(newProject.Id))); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userid, err := strconv.ParseUint(req.UserId, 10, 32)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	newAccess := &models.Access{
		UserId:    uint(userid),
		ProjectId: newProject.Id,
	}

	if err := database.DB.Create(newAccess).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"message": "Copied project",
		"data": newProject,
	})
}


