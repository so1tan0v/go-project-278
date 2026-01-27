package postgres

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

/*Открытие соединения с базой данных*/
func Open(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

