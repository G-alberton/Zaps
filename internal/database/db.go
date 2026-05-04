package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq"
)

// faz a conexão com o banco de dados, lembrar de conectar aqui
func Connect() *sql.DB {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL não configurado")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Erro ao Conectar:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Banco não Reponde", err)
	}

	log.Println("Conectado ao banco")

	return db
}
