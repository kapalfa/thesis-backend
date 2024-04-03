package config

import (
	"context"
	"log"
	"os"

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
	bucketName := os.Getenv("BUCKET_NAME")
	Bucket = client.Bucket(bucketName) //create a bucket instance
}
