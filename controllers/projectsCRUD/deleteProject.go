package projectsCRUD

import (
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
	"github.com/kapalfa/go/database"
	"github.com/kapalfa/go/models"
	"google.golang.org/api/iterator"
)

func DeleteProject(w http.ResponseWriter, r *http.Request) {
	ctx := config.Ctx
	bkt := config.Bucket
	vars := mux.Vars(r)
	id := vars["id"]

	if err := database.DB.Where("project_id = ?", id).Delete(&models.Access{}).Error; err != nil { // delete all access to this project
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := database.DB.Where("id = ?", id).Delete(&models.Project{}).Error; err != nil { // delete project
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	dirName := id + "/"
	it := bkt.Objects(ctx, &storage.Query{Prefix: dirName})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		obj := bkt.Object(attrs.Name)
		if err := obj.Delete(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	http.Error(w, "Deleted project", http.StatusOK)
}
