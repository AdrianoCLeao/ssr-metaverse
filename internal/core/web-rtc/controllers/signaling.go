package controllers

import (
	"github.com/gin-gonic/gin"
	"ssr-metaverse/internal/core/web-rtc/services" 
	"ssr-metaverse/internal/utils"
	"ssr-metaverse/internal/core/error"
)

var chatHub = services.NewHub()

func init() {
	go chatHub.Run()
}

func ChatHandler(c *gin.Context) {
	conn, err := utils.Upgrade(c.Writer, c.Request)
	if err != nil {
		error.RespondWithError(c, *err)
		return
	}

	client := &services.Client{
		Hub:  chatHub,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	chatHub.Register <- client

	go client.WritePump()
	client.ReadPump()
}
