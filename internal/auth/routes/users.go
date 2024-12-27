package routes

import (
	"platable-backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine) {
	users := r.Group("/users")
	{
		users.POST("/", controllers.CreateUser)
		users.GET("/:id", controllers.GetUser)
		users.PUT("/:id", controllers.UpdateUser)
		users.DELETE("/:id", controllers.DeleteUser)
	}
}