package controllers

import (
	"log"
	"net/http"
	"ssr-metaverse/internal/core/objects/services"

	"github.com/gin-gonic/gin"
)

type ObjectController struct {
	Service *services.ObjectService
}

func NewObjectController(service *services.ObjectService) *ObjectController {
	return &ObjectController{Service: service}
}

func (c *ObjectController) UploadObject(ctx *gin.Context) {
	log.Println("⚙️  [Controller] Iniciando UploadObject...")

	bucketName := ctx.PostForm("bucket")
	objectName := ctx.PostForm("object")
	log.Printf("📦 Bucket: %s | 🧱 Object: %s\n", bucketName, objectName)

	file, err := ctx.FormFile("file")
	if err != nil {
		log.Println("❌ [Controller] Erro ao pegar arquivo:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	log.Printf("📁 [Controller] Recebeu arquivo: %s (%d bytes)\n", file.Filename, file.Size)

	metadata := map[string]string{
		"content_type": ctx.PostForm("content_type"),
		"author":       ctx.PostForm("author"),
		"description":  ctx.PostForm("description"),
		"version":      ctx.PostForm("version"),
	}
	log.Printf("📝 [Controller] Metadados recebidos: %+v\n", metadata)

	err = c.Service.UploadObject(bucketName, objectName, file, metadata)
	if err != nil {
		log.Println("❌ [Controller] Erro no serviço:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	log.Println("✅ [Controller] Upload finalizado com sucesso.")
	ctx.JSON(http.StatusOK, gin.H{"message": "Object uploaded successfully"})
}

func (c *ObjectController) ListObjects(ctx *gin.Context) {
	bucketName := ctx.Param("bucket")

	objects, err := c.Service.ListObjects(bucketName)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"objects": objects})
}
