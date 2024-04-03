package filesCRUD

import (
	"encoding/json"
	"net/http"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/gorilla/mux"
	"github.com/kapalfa/go/config"
	"google.golang.org/api/iterator"
)

type Dir struct {
	Name     string          `json:"name"`
	IsDir    bool            `json:"isDir"`
	Filepath string          `json:"filepath,omitempty"`
	Children map[string]*Dir `json:"children"`
}

func GetFiles(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	dirPath := id + "/"
	root := &Dir{Name: id, IsDir: true, Children: make(map[string]*Dir), Filepath: dirPath}

	ctx := config.Ctx
	bkt := config.Bucket
	query := &storage.Query{Prefix: dirPath}
	it := bkt.Objects(ctx, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		relPath := strings.TrimPrefix(attrs.Name, dirPath)
		parts := strings.Split(relPath, string('/'))
		currentDir := root
		for i, part := range parts {
			if part == "" {
				continue
			}
			isDir := i != len(parts)-1
			if _, ok := currentDir.Children[part]; !ok {
				newDir := &Dir{Name: part, IsDir: isDir, Children: make(map[string]*Dir)}
				newDir.Filepath = attrs.Name
				currentDir.Children[part] = newDir
			}
			currentDir = currentDir.Children[part]
		}
	}
	rootJson, err := json.Marshal(root)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(rootJson)
}
