package filesCRUD

import (
	"fmt"
	"io"
	"net/http"

	"github.com/kapalfa/go/config"
)

func SaveFile(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[len("/saveFile/"):]
	ctx := config.Ctx
	bkt := config.Bucket
	obj := bkt.Object(path)

	file, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File: ", file)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	writer := obj.NewWriter(ctx)
	if _, err := io.Copy(writer, file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := writer.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
