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
	filePath := ctx.PostForm("path")
	contentType := ctx.PostForm("content_type")

	err := c.Service.UploadObject(bucketName, objectName, filePath, contentType)
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
