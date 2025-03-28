package web_rtc

import (
	"github.com/gin-gonic/gin"
	"ssr-metaverse/internal/core/web-rtc/controllers" 
)

func RegisterWebRTCRoutes(router *gin.Engine) {
	group := router.Group("/webrtc")
	{
		group.GET("/chat", controllers.ChatHandler)
		group.GET("/ws", controllers.MediaHandler)
	}
}
