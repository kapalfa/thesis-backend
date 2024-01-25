package filesCRUD

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]

	if _, err := os.Stat(path); os.IsNotExist(err) { //kanonika den prepei na yparxei to path idi
		os.MkdirAll(path, 0755)
	}

	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := handler.Filename
	targetPath := filepath.Join(path, filename)
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		response := map[string]interface{}{
			"status": http.StatusBadRequest,
			"message": "File already exists",
		}
		json.NewEncoder(w).Encode(response)
		return
	}
	targetFile, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, file); err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File uploaded successfully"))
}