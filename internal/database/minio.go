package database

import (
	"context"
	"fmt"
	"log"
	"ssr-metaverse/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioInterface interface {
	Connect() error
	BucketExists(bucketName string) (bool, error)
	CreateBucket(bucketName string) error
	UploadObject(bucketName, objectName, filePath, contentType string) error
}

type MinioStorage struct {
	Client *minio.Client
}

var MinioInstance MinioInterface

func (m *MinioStorage) Connect() error {
	endpoint := config.MinioEndpoint
	accessKey := config.MinioAccessKey
	secretKey := config.MinioSecretKey

	log.Println("%w", endpoint)

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return fmt.Errorf("error connecting to MinIO: %w", err)
	}

	m.Client = client
	log.Println("MinIO connected successfully!")
	return nil
}

func (m *MinioStorage) BucketExists(bucketName string) (bool, error) {
	ctx := context.Background()
	exists, err := m.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return false, fmt.Errorf("error checking bucket %s existence: %w", bucketName, err)
	}
	return exists, nil
}

func (m *MinioStorage) CreateBucket(bucketName string) error {
	ctx := context.Background()
	exists, err := m.BucketExists(bucketName)
	if err != nil {
		return err
	}
	if exists {
		log.Printf("Bucket %s already exists\n", bucketName)
		return nil
	}

	err = m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	if err != nil {
		return fmt.Errorf("error creating bucket %s: %w", bucketName, err)
	}

	log.Printf("Bucket %s created successfully\n", bucketName)
	return nil
}

func (m *MinioStorage) UploadObject(bucketName, objectName, filePath, contentType string) error {
	ctx := context.Background()

	_, err := m.Client.FPutObject(ctx, bucketName, objectName, filePath, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("error uploading object %s: %w", objectName, err)
	}

	log.Printf("Object %s uploaded to bucket %s\n", objectName, bucketName)
	return nil
}
