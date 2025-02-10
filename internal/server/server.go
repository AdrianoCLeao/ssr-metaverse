package server

import (
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	"ssr-metaverse/internal/core/auth/routes"
	"ssr-metaverse/internal/middlewares"
	"ssr-metaverse/internal/core/web-rtc/routes" 
)

// HelloHandler godoc
// @Summary Returns a message from the API
// @Description Returns a Hello World message from the API
// @Tags hello
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string "Hello Message"
// @Router /hello [get]
func HelloHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Hello from RESTful API!"})
}

func (s *Server) Start(addr string) error {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	router.Use(middlewares.ErrorHandler())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	routes.RegisterAuthRoutes(router, s.DB)
	routes.RegisterUserRoutes(router, s.DB)
	routes.RegisterProtectedRoutes(router)
	web_rtc.RegisterWebRTCRoutes(router)

	router.StaticFS("/assets", http.Dir("./assets"))

	router.GET("/hello", HelloHandler)

	router.GET("/ws", func(c *gin.Context) {
		s.HandleWebSocket(c.Writer, c.Request)
	})
	

	router.GET("/health", func(c *gin.Context) {
		err := s.DB.CheckHealth()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":  "unhealthy",
				"message": "Error to connect to database",
				"error":   err.Error(),
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "Database is healthy and ready!",
		})
	})

	log.Printf("Starting server in %s...", addr)
	return router.Run(addr)
}
