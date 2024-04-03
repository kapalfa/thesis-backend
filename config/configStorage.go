package config

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
)

var Bucket *storage.BucketHandle
var Ctx context.Context

func ConfigStorage() {
	Ctx = context.Background()
	client, err := storage.NewClient(Ctx)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()
	bucketName := "bucket-editor-files-1312"
	Bucket = client.Bucket(bucketName) //create a bucket instance
}
