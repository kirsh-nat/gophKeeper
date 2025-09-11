package storage

import (
	"context"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func InitMinio() (*minio.Client, error) {
	endpoint := os.Getenv("MINIO_ENDPOINT")    // localhost:9000
	accessKey := os.Getenv("MINIO_ACCESS_KEY") // minioadmin
	secretKey := os.Getenv("MINIO_SECRET_KEY") // minioadmin

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func UploadObject(client *minio.Client, bucket, objectName string, reader *os.File, size int64, contentType string) error {
	_, err := client.PutObject(
		context.Background(),
		bucket,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{ContentType: contentType},
	)
	return err
}

func DownloadObject(client *minio.Client, bucket, objectName string) (*minio.Object, error) {
	return client.GetObject(context.Background(), bucket, objectName, minio.GetObjectOptions{})
}
