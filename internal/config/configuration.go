package config

import "os"

var (
	JwtSecret  = []byte(os.Getenv("JWT_SECRET"))

	DBHost     = os.Getenv("DB_HOST")
	DBPort     = os.Getenv("DB_PORT")
	DBUser     = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName     = os.Getenv("DB_NAME")

	MinioEndpoint = os.Getenv("MINIO_ENDPOINT")
	MinioAccessKey = os.Getenv("MINIO_ACCESSKEY")
	MinioSecretKey = os.Getenv("MINIO_SECRET")

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")

	MongoURI = os.Getenv("MONGO_HOST")
	MongoDBName = os.Getenv("MONGO_NAME")
)
