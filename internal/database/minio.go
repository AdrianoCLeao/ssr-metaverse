package database

import (
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"ssr-metaverse/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioInterface interface {
	Connect() error
	BucketExists(bucketName string) (bool, error)
	CreateBucket(bucketName string) error
	UploadObjectFromReader(bucketName, objectName string, reader multipart.File, size int64, contentType string, metadata map[string]string) error
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

func (m *MinioStorage) UploadObjectFromReader(bucketName, objectName string, reader multipart.File, size int64, contentType string, metadata map[string]string) error {
	ctx := context.Background()

	exists, err := m.Client.BucketExists(ctx, bucketName)
	if err != nil {
		return fmt.Errorf("error checking if bucket exists: %w", err)
	}

	if !exists {
		err = m.Client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}
	}

	_, err = m.Client.PutObject(
		ctx,
		bucketName,
		objectName,
		reader,
		size,
		minio.PutObjectOptions{
			ContentType:  contentType,
			UserMetadata: metadata,
		},
	)
	if err != nil {
		return fmt.Errorf("error uploading object to MinIO: %w", err)
	}
	return nil
}


