package filesCRUD

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
)

type RequestFile struct {
	Name string `json:"name"`
}

func CreateFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	var file RequestFile
	err := json.NewDecoder(r.Body).Decode(&file)
	if err != nil {
		http.Error(w, "Can't create file", http.StatusBadRequest)
		return
	}

	targetPath := path + file.Name
	bkt := config.Bucket
	ctx := config.Ctx
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
	if err := wc.Close(); err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File created successfully"))
}
