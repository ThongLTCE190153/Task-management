package db

import (
	"database/sql"
	"log"

	"trithong.com/task-golang/internal/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectPostgres(cfg config.Config) *sql.DB {
	database, err := sql.Open("pgx", cfg.DatabaseDSN())
	if err != nil {
		log.Fatal("Failed to open database: ", err)
	}

	err = database.Ping()
	if err != nil {
		log.Fatal("Failed to connect database: ", err)
	}

	log.Println("Connected to PostgreSQL successfully")
	return database
}
