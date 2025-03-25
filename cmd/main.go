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
	database.MinioInstance = minio

	mongo := &database.Mongo{}
	if err := mongo.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao MongoDB: %v", err)
	}
	database.MongoInstance = mongo

	redis := &database.Redis{}
	if err := redis.Connect(); err != nil {
		log.Fatalf("Erro ao conectar ao Redis: %v", err)
	}
	database.RedisInstance = redis

	swagger.SwaggerInfo.Title = "SSR Metaverse API"
	swagger.SwaggerInfo.Description = "This is an example API to SSR Metaverse."
	swagger.SwaggerInfo.Version = "1.0"
	swagger.SwaggerInfo.Host = "localhost:8080"
	swagger.SwaggerInfo.BasePath = "/"
	swagger.SwaggerInfo.Schemes = []string{"http"}

	srv := server.NewServer(db, minio, mongo, redis)

	log.Println("Starting server on :8080")

	if err := srv.Start(":8080"); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
