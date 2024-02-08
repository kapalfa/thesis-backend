package filesCRUD

import (
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"os"
)
type Folder struct {
	FolderName string `json:"foldername"`
}
func CreateFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["folderpath"]
	var folder Folder
	json.NewDecoder(r.Body).Decode(&folder)

	targetPath := path + "/" + folder.FolderName 

	log.Println("targetPath ", targetPath)
	if _, err := os.Stat(targetPath); !os.IsNotExist(err) {
	 	response := map[string]interface{} {
	 		"status": http.StatusBadRequest,
	 		"message":"Folder already exists",
	 	}
	 	json.NewEncoder(w).Encode(response)
	 	return 
	} else {
	 	err := os.Mkdir(targetPath, 0755)
	 	if err != nil {
	 		http.Error(w, "Error creating foler", http.StatusInternalServerError)
	 		return 
	 	}
	 	w.WriteHeader(http.StatusOK)
	 	w.Write([]byte("Folder created successfully"))
	}
	
}