package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ssr-metaverse/internal/core/web-rtc/services" 
	"ssr-metaverse/internal/utils"
)

var hub = services.NewHub()

func init() {
	go hub.Run() 
}

func SignalWs(c *gin.Context) {
	conn, err := utils.Upgrade(c.Writer, c.Request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error establishing websocket connection."})
		return
	}

	client := &services.Client{
		Hub:  hub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}
	client.Hub.Register <- client

	go client.WritePump()
	client.ReadPump() 
}
