package routes

import (
	"ssr-metaverse/internal/core/auth/controllers"
	"ssr-metaverse/internal/core/auth/services"
	"ssr-metaverse/internal/database"

	"github.com/gin-gonic/gin"
)

func RegisterUserRoutes(r *gin.Engine, db database.DBInterface) {
	userService := &services.UserService{DB: db}
	userController := controllers.NewUserController(userService)

	users := r.Group("/users")
	{
		users.POST("/", userController.CreateUser)
		users.GET("/:id", userController.GetUser)
		users.PUT("/:id", userController.UpdateUser)
		users.DELETE("/:id", userController.DeleteUser)
	}
}