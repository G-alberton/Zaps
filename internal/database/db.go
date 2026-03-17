package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// faz a conexão com o banco de dados, lembrar de conectar aqui
func Conect() *sql.DB {
	connStr := "Acesso ao Banco de Dados"

	db, err := sql.Open("Postgres", connStr)
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

/*
	tabela tem que ser:
	create table contacts (
		phone TEXT Primary Key,
		name TEXT
	);
*/
