package server

import (
	"ssr-metaverse/internal/auth/routes"
	"ssr-metaverse/internal/database"

	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

/*
Starts the server using gin for HTTP and WebSocket endpoints.
It configures CORS, serves static files from './assets' and sets up the routes.
*/
func (s *Server) Start(addr string) error {
	router := gin.Default()
	database.Connect()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.RegisterAuthRoutes(router)
	routes.RegisterUserRoutes(router)
	routes.RegisterProtectedRoutes(router)

	router.StaticFS("/assets", http.Dir("./assets"))

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from RESTful API!"})
	})

	router.GET("/ws", func(c *gin.Context) {
		s.HandleWebSocket(c.Writer, c.Request)
	})

	router.GET("/health", func(c *gin.Context) {
		err := database.CheckHealth()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "unhealthy",
				"message": "Erro ao conectar ao banco de dados",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Banco de dados conectado com sucesso",
		})
	})

	return router.Run(addr)
}
