package main

import (
	"log"
	"ssr-metaverse/docs"
	_ "ssr-metaverse/docs"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/server"
)

func main() {
	db := &database.Database{}
	if err := db.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	docs.SwaggerInfo.Title = "SSR Metaverse API"
    docs.SwaggerInfo.Description = "This is an example API to SSR Metaverse."
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "localhost:8080"
    docs.SwaggerInfo.BasePath = "/"
    docs.SwaggerInfo.Schemes = []string{"http"}

	srv := server.NewServer(db)
	
	log.Println("Starting server on :8080")
	
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
