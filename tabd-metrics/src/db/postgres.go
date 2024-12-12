package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func InitPostgres() (*sql.DB, error) {
	connStr := "postgres://root:root@localhost:5432/root?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("erro ao conectar ao PostgreSQL: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(50),
		email VARCHAR(50),
		preferences TEXT
	)`)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar tabela: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS books (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(50),
		author VARCHAR(50),
		genre VARCHAR(50),
		description TEXT
	)`)

	if err != nil {
		return nil, fmt.Errorf("erro ao criar tabela: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS ratings (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		user_id UUID,
		book_id UUID,
		note INT,
		comment TEXT
	)`)

	return db, nil
}
