package filesCRUD

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"os"
)

func CreateFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	filename := struct {	
		Filename string `json:"filename"`
	}{}
	json.NewDecoder(r.Body).Decode(&filename)
	targetPath := path + "/" + filename.Filename

	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
		response := map[string]interface{}{
			"status": http.StatusBadRequest,
			"message": "File already exists",
		}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		file, err := os.Create(targetPath)
		if err != nil {
			http.Error(w, "Error creating file", http.StatusInternalServerError)
			return 
		}
		defer file.Close()
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File created successfully"))
	}
}