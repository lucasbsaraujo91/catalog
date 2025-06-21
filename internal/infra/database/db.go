package database

import (
	"catalog/configs"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(cfg *configs.Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	db, err := sql.Open(cfg.DBDriver, dsn)
	if err != nil {
		log.Fatalf("Erro ao conectar no banco: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Erro ao pingar no banco: %v", err)
	}

	log.Println("Conexão com banco de dados realizada com sucesso.")
	return db
}
