package objects

import (
	"ssr-metaverse/internal/core/objects/controllers"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/core/objects/services"

	"github.com/gin-gonic/gin"
)

func RegisterObjectRoutes(router *gin.Engine) {
	minioService := services.NewObjectService(database.MinioInstance)
	objectController := controllers.NewObjectController(minioService)

	objects := router.Group("/objects")
	{
		objects.POST("/upload", objectController.UploadObject)
		objects.GET("/list/:bucket", objectController.ListObjects)
	}
}
