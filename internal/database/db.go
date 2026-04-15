package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// faz a conexão com o banco de dados, lembrar de conectar aqui
func Connect() *sql.DB {
	connStr := "Acesso ao Banco de Dados"

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
