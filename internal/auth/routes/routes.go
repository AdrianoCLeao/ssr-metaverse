package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	RegisterAuthRoutes(r)
	RegisterUserRoutes(r)

	RegisterProtectedRoutes(r)

	return r
}