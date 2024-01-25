package filesCRUD

import (
	"encoding/json"
	"io"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"path/filepath"
	"strings"

)

func UploadFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	var filename string

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}
	files := r.MultipartForm.File["folder"]
	for i, file := range files {
		contentDisposition := file.Header.Get("Content-Disposition")
		parts := strings.Split(contentDisposition, ";")
		for _, part := range parts {
			if strings.Contains(part, "filename") {
				filename = strings.Split(part, "=")[1]
				filename = strings.Trim(filename, "\"")
				if(i == 0){
					tmp := strings.Split(filename, "/")
					outermostDir := tmp[0]
					if _,err := os.Stat(path + "/" + outermostDir); !os.IsNotExist(err) {
						response := map[string]interface{}{
							"status": http.StatusBadRequest,
							"message": "Folder already exists",
						}
						json.NewEncoder(w).Encode(response)
						return
					} 
				}
			}
		}
		outFilePath := path + "/" + filename
		dir := filepath.Dir(outFilePath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			os.MkdirAll(dir, 0755)
		} 
		if _, err := os.Stat(outFilePath); !os.IsNotExist(err) {
			response := map[string]interface{}{
				"status": http.StatusBadRequest,
				"message": "File already exists",
			}
			json.NewEncoder(w).Encode(response)
			return
		}

	  	fileContent, err := file.Open()
	  	if err != nil {
	  		http.Error(w, "Can't read file", http.StatusBadRequest)
			return
	  	}
	 	defer fileContent.Close()
		outFile, err := os.Create(outFilePath)
		if err != nil {
			http.Error(w, "Error on folder upload", http.StatusBadRequest)
			return
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, fileContent)
		if err != nil {
			http.Error(w, "Error on folder upload", http.StatusBadRequest)
			return
		}
	}
	w.Write([]byte("Folder uploaded"))
}