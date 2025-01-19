package routes

import (
	"ssr-metaverse/internal/middlewares"
	
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	RegisterAuthRoutes(r)
	RegisterUserRoutes(r)

	protected := r.Group("/protected")
	protected.Use(middlewares.JWTMiddleware())
	{
		protected.GET("/profile", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "This is a protected route, accessible only with a valid token!"})
		})
	}

	return r
}