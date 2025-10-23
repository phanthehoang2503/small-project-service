package main

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDB(dsn string) {
	if dsn == "" {
		dsn = "postgres://postgress:25032004@localhost:5432/postgres?sslmode=disable"
	}

	DB, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect db: ", err)
	}

	DB.SetMaxOpenConns(10)
	DB.SetMaxIdleConns(5)
}

func MigrateProduct(db *sqlx.DB) {
	schema := `
		CREATE TABLE IF NOT EXISTS products (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT,
		price NUMERIC NOT NULL
	);`

	if _, err := db.Exec(schema); err != nil {
		log.Fatal("migrate failed:", err)
	}
}
