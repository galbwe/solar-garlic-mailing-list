package jobs

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// upload a file to s3
func UploadToS3(filename string, bucket string, key string) error {
	// The session the S3 Uploader will use
	sess := session.Must(session.NewSession())

	// Create an uploader with the session and default options
	uploader := s3manager.NewUploader(sess)

	f, err := os.Open(filename)
	if err != nil {
		slog.Error("Failed to open file", "filename", filename, "err", err)
		return fmt.Errorf("failed to open file %q, %v", filename, err)
	}

	// Upload the file to S3.
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   f,
	})
	if err != nil {
		slog.Error("Failed to upload file", "bucket", bucket, "key", key, "err", err)
		return fmt.Errorf("failed to upload file, %v", err)
	}
	slog.Info("file uploaded to s3", "location", aws.StringValue(&result.Location))
	return nil
}
