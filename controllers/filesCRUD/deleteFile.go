package filesCRUD

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
)

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]

	ctx := config.Ctx
	bkt := config.Bucket
	obj := bkt.Object(path)

	if err := obj.Delete(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File deleted successfully"))
}
