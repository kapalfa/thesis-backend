package filesCRUD

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
)
type RequestFolder struct {
	Name string `json:"name"`
}
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["folderpath"]
	var folder RequestFolder
	err := json.NewDecoder(r.Body).Decode(&folder)
	if err != nil {
		http.Error(w, "Can't create folder", http.StatusBadRequest)
		return
	}
	targetPath := path + folder.Name + "/"

	bkt := config.Bucket
	ctx := config.Ctx

	obj := bkt.Object(targetPath)
	_, err = obj.Attrs(ctx) // Check if folder already exists
	if err == nil {
		response := map[string]interface{}{
			"status": http.StatusBadRequest,
			"message": "Folder already exists",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	writer := obj.NewWriter(ctx)
	if err := writer.Close(); err != nil {
		http.Error(w, "Error creating folder", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Folder created successfully"))
}