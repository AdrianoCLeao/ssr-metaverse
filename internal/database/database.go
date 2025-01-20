package database

import (
	"database/sql"
	"fmt"
	"log"
	"ssr-metaverse/internal/config"

	_ "github.com/lib/pq"
)

type DBInterface interface {
	Connect() error
	CheckHealth() error
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type Database struct {
	DB *sql.DB
}

var DBInstance DBInterface

func (d *Database) Connect() error {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost,
		config.DBPort,
		config.DBUser,
		config.DBPassword,
		config.DBName,
	)

	var err error
	d.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("erro ao conectar ao banco de dados: %w", err)
	}

	if err = d.DB.Ping(); err != nil {
		return fmt.Errorf("erro ao verificar a conexão com o banco de dados: %w", err)
	}

	log.Println("Banco de dados conectado com sucesso!")
	return nil
}

func (d *Database) CheckHealth() error {
	err := d.DB.Ping()
	if err != nil {
		log.Printf("Erro ao verificar a saúde do banco de dados: %v", err)
	}
	return err
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}
