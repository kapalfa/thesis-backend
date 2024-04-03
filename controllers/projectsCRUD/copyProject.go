package projectsCRUD

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"google.golang.org/api/iterator"
)

type CopyRequest struct {
	ProjectId string `json:"projectid"`
	UserId    string `json:"userid"`
}

func copyDir(ctx context.Context, bucket *storage.BucketHandle, srcPrefix, dstPrefix string) error {
	it := bucket.Objects(ctx, &storage.Query{Prefix: srcPrefix})

	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		srcObject := bucket.Object(attrs.Name)
		dstObject := bucket.Object(strings.Replace(attrs.Name, srcPrefix, dstPrefix, 1))

		if _, err := dstObject.CopierFrom(srcObject).Run(ctx); err != nil {
			return err
		}
	}
	return nil
}
func CopyProject(w http.ResponseWriter, r *http.Request) {
	ctx := config.Ctx
	bkt := config.Bucket
	var req CopyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	project := &models.Project{}
	if err := database.DB.Where("id = ? AND public = ?", req.ProjectId, true).First(project).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newProject := &models.Project{
		Name:        project.Name,
		Description: project.Description,
		Public:      false,
	}

	if err := database.DB.Create(newProject).Error; err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	srcPrefix := req.ProjectId + "/"
	dstPrefix := strconv.FormatUint(uint64(newProject.Id), 10) + "/"
	if err := copyDir(ctx, bkt, srcPrefix, dstPrefix); err != nil {
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
		"status":  "success",
		"message": "Copied project",
		"data":    newProject,
	})
}
