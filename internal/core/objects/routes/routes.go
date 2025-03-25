package objects

import (
	"ssr-metaverse/internal/core/objects/controllers"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/core/objects/services"

	"github.com/gin-gonic/gin"
)

func RegisterObjectRoutes(router *gin.Engine, storage database.MinioInterface, mongo database.MongoInterface, redis database.RedisInterface) {
	objectService := services.NewObjectService(storage, mongo, redis)
	objectController := controllers.NewObjectController(objectService)

	objects := router.Group("/objects")
	{
		objects.POST("/upload", objectController.UploadObject)
		objects.GET("/list/:bucket", objectController.ListObjects)
	}
}

