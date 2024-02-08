package filesCRUD

import (
	"net/http"
	"os"
	"github.com/gorilla/mux"
)

func GetFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	if _, err := os.Stat(path); os.IsNotExist(err) {
		http.Error(w, "File does not exist", http.StatusBadRequest)
		return
	}

	http.ServeFile(w, r, path)
}