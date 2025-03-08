package database

import (
	"context"
	"fmt"
	"log"
	"ssr-metaverse/internal/config"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisInterface interface {
	Connect() error
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Ping() error
}

type Redis struct {
	Client *redis.Client
}

var RedisInstance RedisInterface

func (r *Redis) Connect() error {
	r.Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort),
		Password: config.RedisPassword,
		DB:       0,
	})

	if err := r.Ping(); err != nil {
		return fmt.Errorf("erro ao conectar ao Redis: %w", err)
	}

	log.Println("Redis conectado com sucesso!")
	return nil
}

func (r *Redis) Set(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()
	return r.Client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(key string) (string, error) {
	ctx := context.Background()
	return r.Client.Get(ctx, key).Result()
}

func (r *Redis) Del(key string) error {
	ctx := context.Background()
	return r.Client.Del(ctx, key).Err()
}

func (r *Redis) Ping() error {
	ctx := context.Background()
	_, err := r.Client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Erro ao verificar a conexão com o Redis: %v", err)
	}
	return err
}
