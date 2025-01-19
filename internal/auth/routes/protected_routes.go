package routes

import (
	"ssr-metaverse/internal/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterProtectedRoutes(router *gin.Engine) {
	protected := router.Group("/protected")

	protected.GET("/admin", middlewares.RolesGuard("admin"), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Welcome, Admin!"})
	})

	protected.GET("/profile", middlewares.JWTMiddleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "This is a protected route, accessible only with a valid token!"})
	})

}
