package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"ssr-metaverse/internal/config"
)

var DB *sql.DB

func Connect() {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	var err error
	DB, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco de dados: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Erro ao verificar a conexão com o banco de dados: %v", err)
	}

	log.Println("Banco de dados conectado com sucesso!")
}

func CheckHealth() error {
    err := DB.Ping()
    if err != nil {
        log.Printf("Erro ao verificar a saúde do banco de dados: %v", err)
    }
    return err
}
