package routes

import (
	"ssr-metaverse/internal/auth/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
	}
}