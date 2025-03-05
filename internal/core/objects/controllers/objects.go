package controllers

import (
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
	bucketName := ctx.PostForm("bucket")
	objectName := ctx.PostForm("object")

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	metadata := map[string]string{
		"content_type": ctx.PostForm("content_type"),
		"author":       ctx.PostForm("author"),
		"description":  ctx.PostForm("description"),
		"version":      ctx.PostForm("version"),
	}

	err = c.Service.UploadObject(bucketName, objectName, file, metadata)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
