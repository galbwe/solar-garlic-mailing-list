package jobs

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// upload a file to s3
func UploadToS3(filename string, bucket string, key string) error {
	// The session the S3 Uploader will use
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        log.Fatal(err)
    }

	// Create an uploader with the session and default options
	client := s3.NewFromConfig(cfg)

	f, err := os.Open(filename)
	if err != nil {
		slog.Error("Failed to open file", "filename", filename, "err", err)
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}
	defer f.Close()

	// Upload the file to S3.
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		slog.Error("Failed to upload file", "bucket", bucket, "key", key, "err", err)
		return fmt.Errorf("failed to upload file, %v", err)
	}
	slog.Info("file uploaded to s3", "bucket", bucket, "key", key)
	return nil
}
