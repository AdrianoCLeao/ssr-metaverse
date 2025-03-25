package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mime/multipart"
	"ssr-metaverse/internal/database"
	"time"

	"github.com/minio/minio-go/v7"
	"go.mongodb.org/mongo-driver/bson"
)

type ObjectService struct {
	Minio database.MinioInterface
	MongoDB database.MongoInterface
	Redis   *database.Redis
}

func NewObjectService(minio database.MinioInterface, mongo database.MongoInterface) *ObjectService {
	return &ObjectService{
		Minio: minio,
		MongoDB: mongo,
	}
}

func (s *ObjectService) UploadObject(bucketName, objectName string, file *multipart.FileHeader, metadata map[string]string) error {
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	if s.Minio == nil {
		return fmt.Errorf("MinIO storage was not initialized")
	}

	err = s.Minio.UploadObjectFromReader(bucketName, objectName, src, file.Size, metadata["content_type"], metadata)
	if err != nil {
		return fmt.Errorf("failed to upload object: %w", err)
	}

	metadataDoc := bson.M{
		"bucket":       bucketName,
		"object_name":  objectName,
		"content_type": metadata["content_type"],
		"author":       metadata["author"],
		"description":  metadata["description"],
		"version":      metadata["version"],
		"uploaded_at":  time.Now(),
	}

	_, err = s.MongoDB.InsertOne("objects_metadata", metadataDoc)
	if err != nil {
		return fmt.Errorf("failed to save metadata in MongoDB: %w", err)
	}

	docJSON, err := json.Marshal(metadataDoc)
	if err != nil {
		return fmt.Errorf("failed to serialize metadata for Redis: %w", err)
	}

	err = s.Redis.Set("world-object:"+objectName, docJSON, 3600)
	if err != nil {
		return fmt.Errorf("failed to cache metadata in Redis: %w", err)
	}

	return nil
}


func (s *ObjectService) ListObjects(bucketName string) ([]minio.ObjectInfo, error) {
	ctx := context.Background()
	objectCh := s.Minio.(*database.MinioStorage).Client.ListObjects(ctx, bucketName, minio.ListObjectsOptions{Recursive: true})

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
