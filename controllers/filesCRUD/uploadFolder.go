package filesCRUD

import (
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
)

func UploadFolder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars["filepath"]
	var filename string

	ctx := config.Ctx
	bkt := config.Bucket

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		return
	}
	files := r.MultipartForm.File["folder"]

	for _, file := range files {
		contentDisposition := file.Header.Get("Content-Disposition")
		parts := strings.Split(contentDisposition, ";")
		for _, part := range parts {
			if strings.Contains(part, "filename") {
				filename = strings.Split(part, "=")[1]
				filename = strings.Trim(filename, "\"")
			}
		}
		outFilePath := path + filename
		dir := filepath.Dir(outFilePath)
		if dir != "." {
			folderObj := bkt.Object(dir + "/")
			folderWriter := folderObj.NewWriter(ctx)
			folderWriter.Close()
		}
		obj := bkt.Object(outFilePath)
		writer := obj.NewWriter(ctx)
		_, err := obj.Attrs(ctx)
		if err == nil {
			http.Error(w, "Folder already exists", http.StatusBadRequest)
			return
		}
		fileContent, err := file.Open()
		if err != nil {
			http.Error(w, "Can't read file", http.StatusBadRequest)
			return
		}
		defer fileContent.Close()
		_, err = io.Copy(writer, fileContent)
		if err != nil {
			http.Error(w, "Error on folder upload", http.StatusBadRequest)
			return
		}
		if err = writer.Close(); err != nil {
			http.Error(w, "Error on folder upload", http.StatusBadRequest)
			return
		}
	}
	w.Write([]byte("Folder uploaded"))
}
