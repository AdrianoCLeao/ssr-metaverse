package main

import (
	"log"
	"ssr-metaverse/api/swagger"
	_ "ssr-metaverse/api/swagger"
	"ssr-metaverse/internal/database"
	"ssr-metaverse/internal/server"
)

func main() {
	db := &database.Database{}
	if err := db.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	minio := &database.MinioStorage{}
	if err := minio.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao minio: %v", err)
	}

	swagger.SwaggerInfo.Title = "SSR Metaverse API"
    swagger.SwaggerInfo.Description = "This is an example API to SSR Metaverse."
    swagger.SwaggerInfo.Version = "1.0"
    swagger.SwaggerInfo.Host = "localhost:8080"
    swagger.SwaggerInfo.BasePath = "/"
    swagger.SwaggerInfo.Schemes = []string{"http"}

	srv := server.NewServer(db)
	
	log.Println("Starting server on :8080")
	
	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
