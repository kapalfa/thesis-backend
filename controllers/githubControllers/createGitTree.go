package githubControllers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"strings"

	"cloud.google.com/go/storage"
	"github.com/google/go-github/github"
	"github.com/kapalfa/go/config"
	"golang.org/x/oauth2"
	"google.golang.org/api/iterator"
)

func CreateGitTreeOld(accessToken, username, repo, branch, projectFolder, sha string) {
	ctx := config.Ctx
	bkt := config.Bucket
	context := context.Background()
	query := &storage.Query{Prefix: projectFolder}
	it := bkt.Objects(ctx, query)
	ghClient := github.NewClient(nil)
	//tree := make(map[string]interface{})
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		r, err := bkt.Object(attrs.Name).NewReader(ctx)
		if err != nil {
			log.Printf("Error reading object: %v", err)
			return
		}
		content, _ := io.ReadAll(r)
		r.Close()
		log.Println("content: ", string(content))
		log.Print("attrs.Name: ", attrs.Name)

		blob, _, _ := ghClient.Git.CreateBlob(context, username, repo, &github.Blob{
			Content:  github.String(string(content)),
			Encoding: github.String("utf-8"),
		})

		log.Println(blob)
	}
}
func IsDir(obj *storage.ObjectAttrs) bool {
	return strings.HasSuffix(obj.Name, "/")
}
func CreateGitTreeMiddle(accessToken, username, repo, branch, projectFolder, sha string) {
	ctx := config.Ctx
	bkt := config.Bucket
	context := context.Background()
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	ghClient := github.NewClient(oauthClient)

	objIter := bkt.Objects(ctx, &storage.Query{Prefix: projectFolder})
	var objects []*storage.ObjectAttrs
	for {
		attrs, err := objIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error reading object: %v", err)
			return
		}
		objects = append(objects, attrs)
	}
	//tree := make(map[string]interface{})
	var tree []github.TreeEntry
	for _, obj := range objects {
		filename := strings.TrimPrefix(obj.Name, projectFolder)

		if IsDir(obj) {
			path := strings.TrimSuffix(filename, "/")
			if path != "" {
				tree = append(tree, github.TreeEntry{
					Path: &path,
					Mode: github.String("040000"),
					Type: github.String("tree"),
				})
			}
			//	tree[path] = map[string]interface{}{"type": "tree"}
		} else {
			path := filename
			obj := bkt.Object(obj.Name)
			data, err := obj.NewReader(ctx)
			if err != nil {
				log.Printf("Error reading object: %v", err)
				return
			}
			defer data.Close()

			content, err := io.ReadAll(data)
			if err != nil {
				log.Printf("Error reading object data: %v", err)
				return
			}
			blob, _, err := ghClient.Git.CreateBlob(context, username, repo, &github.Blob{
				Content:  github.String(string(content)),
				Encoding: github.String("utf-8"),
			})
			if err != nil {
				log.Printf("Error creating blob: %v", err)
				return
			}
			//	sha := fmt.Sprintf("sha1:%x", sha1.Sum(content))
			tree = append(tree, github.TreeEntry{
				SHA:  blob.SHA,
				Path: github.String(path),
				Mode: github.String("100644"),
				Type: github.String("blob"),
			})
			// tree[path] = map[string]interface{}{
			// 	"type": "blob",
			// 	"mode": "100644",
			// 	"sha":  fmt.Sprintf("sha1:%x", sha1.Sum(content)),
			// }
		}
	}
	treeJson, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		log.Printf("Error marshalling newTreeRef: %v", err)
		return
	}
	newTreeRef, _, err := ghClient.Git.CreateTree(context, username, repo, "", tree)
	if err != nil {
		log.Printf("Error creating tree: %v", err)
		return
	}
	log.Print("newTreeRef: ", newTreeRef)
	log.Print("treeJson: ", string(treeJson))
}
func CreateGitTree(ghClient *github.Client, accessToken, username, repo, branch, projectFolder string) (*github.Tree, error) {
	ctx := config.Ctx
	bkt := config.Bucket
	query := &storage.Query{Prefix: projectFolder}
	it := bkt.Objects(ctx, query)
	var objects []*storage.ObjectAttrs
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error reading object: %v", err)
			return nil, err
		}
		objects = append(objects, attrs)
	}
	var treeEntries []github.TreeEntry
	for _, obj := range objects {
		filename := strings.TrimPrefix(obj.Name, projectFolder)

		if IsDir(obj) {
			path := strings.TrimSuffix(filename, "/")
			if path != "" {

				subTree, err := CreateGitTree(ghClient, accessToken, username, repo, branch, obj.Name)
				if err != nil {
					log.Printf("Error creating sub tree: %v", err)
					return nil, err
				}
				treeEntries = append(treeEntries, github.TreeEntry{
					Path: github.String(path),
					Mode: github.String("040000"),
					Type: github.String("tree"),
					SHA:  subTree.SHA,
				})
			}
		} else {
			fileObj := bkt.Object(obj.Name)
			data, err := fileObj.NewReader(ctx)
			if err != nil {
				log.Printf("Error reading object: %v", err)
				return nil, err
			}
			defer data.Close()

			content, err := io.ReadAll(data)
			if err != nil {
				log.Printf("Error reading object data: %v", err)
				return nil, err
			}
			blob, _, err := ghClient.Git.CreateBlob(context.Background(), username, repo, &github.Blob{
				Content:  github.String(string(content)),
				Encoding: github.String("utf-8"),
			})
			if err != nil {
				log.Printf("Error creating blob: %v", err)
				return nil, err
			}
			treeEntries = append(treeEntries, github.TreeEntry{
				Path: github.String(filename),
				Mode: github.String("100644"),
				Type: github.String("blob"),
				SHA:  blob.SHA,
			})
		}
	}
	tree, _, err := ghClient.Git.CreateTree(context.Background(), username, repo, "", treeEntries)
	if err != nil {
		log.Printf("Error creating tree: %v", err)
		return nil, err
	}
	return tree, nil
}
