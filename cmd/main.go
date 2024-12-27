package main

import (
	"log"
	"ssr-metaverse/internal/server"
)

func main() {
	srv := server.NewServer()
	log.Println("Starting server on :8080")
	log.Fatal(srv.Start(":8080"))
}
