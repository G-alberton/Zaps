package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// faz a conexão com o banco de dados, lembrar de conectar aqui
func Connect() *sql.DB {
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
	tabela para o contact:
	create table contacts (
		phone TEXT Primary Key,
		name TEXT
	);
*/

/*
	CREATE TABLE messages (
	id SERIAL PRIMARY KEY,
	from_phone TEXT,
	type TEXT,
	body TEXT,
	media_id TEXT,
	created_at TIMESTAMP DEFAULT NOW()
);

*/
