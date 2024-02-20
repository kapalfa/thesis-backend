package filesCRUD

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]

	ctx := config.Ctx
	bkt := config.Bucket

	//r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	targetPath := path + filename

	obj := bkt.Object(targetPath)
	_, err = obj.Attrs(ctx)
	if err == nil {
		response := map[string]interface{}{
			"message": "File already exists",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	if err := wc.Close(); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}
