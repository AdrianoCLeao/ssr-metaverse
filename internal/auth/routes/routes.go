package routes

import (
	"ssr-metaverse/internal/database"

	"github.com/gin-gonic/gin"
)

func SetupRouter(db database.DBInterface) *gin.Engine {
	r := gin.Default()

	RegisterAuthRoutes(r, db)
	RegisterUserRoutes(r, db)

	RegisterProtectedRoutes(r)

	return r
}