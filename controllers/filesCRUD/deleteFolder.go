package filesCRUD

import (
	"net/http"
	"os"
	"log"
	"github.com/gorilla/mux"
)

func DeleteFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["folderpath"]

	log.Println("Deleting folder: ", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "Folder does not exist", http.StatusBadRequest)
		return
	}

	err := os.RemoveAll(path)
	if err != nil {
		http.Error(w, "Error deleting folder", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Folder deleted successfully"))
}