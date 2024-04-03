package filesCRUD

import (
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	bkt := config.Bucket
	ctx := config.Ctx
	obj := bkt.Object(path) // get an object handler

	reader, err := obj.NewReader(ctx)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	defer reader.Close()

	if _, err := io.Copy(w, reader); err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
}
