package database

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDb(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database Connected Successfully")
	return db, nil
}
