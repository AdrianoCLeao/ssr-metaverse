package server

import (
	"log"
	"net/http"
	"ssr-metaverse/internal/auth/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) Start(addr string) error {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.RegisterAuthRoutes(router, s.DB)
	routes.RegisterUserRoutes(router, s.DB)
	routes.RegisterProtectedRoutes(router)

	router.StaticFS("/assets", http.Dir("./assets"))

	router.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello from RESTful API!"})
	})

	router.GET("/ws", func(c *gin.Context) {
		s.HandleWebSocket(c.Writer, c.Request)
	})

	router.GET("/health", func(c *gin.Context) {
		err := s.DB.CheckHealth()
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

	log.Printf("Iniciando servidor em %s...", addr)
	return router.Run(addr)
}
