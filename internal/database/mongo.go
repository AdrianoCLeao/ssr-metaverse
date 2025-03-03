package database

import (
	"context"
	"fmt"
	"log"
	"ssr-metaverse/internal/config"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoInterface interface {
	Connect() error
	InsertOne(collection string, document interface{}) (*mongo.InsertOneResult, error)
	FindOne(collection string, filter interface{}) (*mongo.SingleResult, error)
	DeleteOne(collection string, filter interface{}) (*mongo.DeleteResult, error)
	Ping() error
}

type Mongo struct {
	Client   *mongo.Client
	Database *mongo.Database
}

var MongoInstance MongoInterface

func (m *Mongo) Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(config.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao MongoDB: %w", err)
	}

	m.Client = client
	m.Database = client.Database(config.MongoDBName)

	if err := m.Ping(); err != nil {
		return fmt.Errorf("erro ao verificar a conexão com o MongoDB: %w", err)
	}

	log.Println("MongoDB conectado com sucesso!")
	return nil
}

func (m *Mongo) InsertOne(collection string, document interface{}) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Database.Collection(collection).InsertOne(ctx, document)
}

func (m *Mongo) FindOne(collection string, filter interface{}) (*mongo.SingleResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	result := m.Database.Collection(collection).FindOne(ctx, filter)
	return result, result.Err()
}

func (m *Mongo) DeleteOne(collection string, filter interface{}) (*mongo.DeleteResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return m.Database.Collection(collection).DeleteOne(ctx, filter)
}

func (m *Mongo) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := m.Client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Erro ao verificar a conexão com o MongoDB: %v", err)
	}
	return err
}
