package services

import (
	"context"
	"fmt"
	"log"
	"ssr-metaverse/internal/database"

	"github.com/minio/minio-go/v7"
)

type ObjectService struct {
	Storage database.MinioInterface
}

func NewObjectService(storage database.MinioInterface) *ObjectService {
	return &ObjectService{Storage: storage}
}

func (s *ObjectService) UploadObject(bucketName, objectName, filePath, contentType string) error {
	err := s.Storage.UploadObject(bucketName, objectName, filePath, contentType)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}
	return nil
}

func (s *ObjectService) ListObjects(bucketName string) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	objectCh := s.Storage.(*database.MinioStorage).Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			log.Printf("Error listing objects: %v", object.Err)
			return nil, object.Err
		}
		objects = append(objects, object)
	}
	return objects, nil
}
