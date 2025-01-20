package routes

import (
	"ssr-metaverse/internal/auth/controllers"
	"ssr-metaverse/internal/auth/services"
	"ssr-metaverse/internal/database"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine, db database.DBInterface) {
	authService := &services.UserService{DB: db}
	authController := controllers.NewAuthController(authService)
	
	auth := r.Group("/auth")
	{
		auth.POST("/login", authController.Login)
	}
}