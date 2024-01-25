package filesCRUD

import (
	"encoding/json"
	"net/http"
	"github.com/gorilla/mux"
	"os"
)

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	
	if _, err := os.Stat(path); os.IsNotExist(err) {
		response := map[string]interface{}{
			"status": http.StatusBadRequest,
			"message": "File does not exist",
		}
		json.NewEncoder(w).Encode(response)
		return
	} else {
		err := os.Remove(path)
		if err != nil {
			http.Error(w, "Error deleting file", http.StatusInternalServerError)
			return 
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("File deleted successfully"))
	}
}

