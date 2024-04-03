package filesCRUD

import (
	"net/http"
	"github.com/gorilla/mux"
	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"github.com/kapalfa/go/config"
)

func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["folderpath"]

	ctx := config.Ctx
	bkt := config.Bucket
	query := &storage.Query{Prefix: path}
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := bkt.Object(attrs.Name).Delete(ctx); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Folder deleted successfully"))
}