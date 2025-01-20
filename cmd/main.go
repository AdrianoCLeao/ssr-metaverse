package main

import (
	"log"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/server"
)

func main() {
	db := &database.Database{}
	if err := db.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	// Criar servidor
	srv := server.NewServer(db)
	
	log.Println("Starting server on :8080")
	
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
