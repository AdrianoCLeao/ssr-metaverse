package routes

import (
	"ssr-metaverse/internal/core/auth/controllers"
	"ssr-metaverse/internal/core/auth/services"
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