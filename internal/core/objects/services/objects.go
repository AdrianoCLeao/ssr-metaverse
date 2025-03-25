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
	Storage database.MinioInterface
	MongoDB *database.Mongo
	Redis   *database.Redis
}

func NewObjectService(storage database.MinioInterface) *ObjectService {
	return &ObjectService{Storage: storage}
}

func (s *ObjectService) UploadObject(bucketName, objectName string, file *multipart.FileHeader, metadata map[string]string) error {
	log.Println("⚙️  [Service] UploadObject iniciado")

	_ = context.Background()
	src, err := file.Open()
	if err != nil {
		log.Println("❌ [Service] Falha ao abrir arquivo:", err)
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()
	log.Printf("📁 [Service] Arquivo aberto: %s (%d bytes)\n", file.Filename, file.Size)

	log.Println("🚀 [Service] Enviando para MinIO...")
	err = s.Storage.UploadObjectFromReader(bucketName, objectName, src, file.Size, metadata["content_type"], metadata)
	if err != nil {
		log.Println("❌ [Service] Falha ao fazer upload no MinIO:", err)
		return fmt.Errorf("failed to upload object: %w", err)
	}

	log.Println("🧾 [Service] Gerando metadados para MongoDB")
	metadataDoc := bson.M{
		"bucket":       bucketName,
		"object_name":  objectName,
		"content_type": metadata["content_type"],
		"author":       metadata["author"],
		"description":  metadata["description"],
		"version":      metadata["version"],
		"uploaded_at":  time.Now(),
	}
	log.Printf("🧾 [Service] Documento: %+v\n", metadataDoc)

	_, err = s.MongoDB.InsertOne("objects_metadata", metadataDoc)
	if err != nil {
		log.Println("❌ [Service] Falha ao salvar no MongoDB:", err)
		return fmt.Errorf("failed to save metadata in MongoDB: %w", err)
	}
	log.Println("✅ [Service] Metadados salvos no MongoDB")

	docJSON, err := json.Marshal(metadataDoc)
	if err != nil {
		log.Println("❌ [Service] Falha ao serializar metadados para Redis:", err)
		return fmt.Errorf("failed to serialize metadata for Redis: %w", err)
	}

	err = s.Redis.Set("world-object:"+objectName, docJSON, 3600)
	if err != nil {
		log.Println("❌ [Service] Falha ao salvar metadados no Redis:", err)
		return fmt.Errorf("failed to cache metadata in Redis: %w", err)
	}
	log.Println("✅ [Service] Metadados salvos no Redis com TTL 3600s")

	log.Println("🎉 [Service] Upload finalizado com sucesso!")
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
