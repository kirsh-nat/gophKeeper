package app

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func InitDB() (*sql.DB, string, error) {
	if err := godotenv.Load("/opt/gophKeeper/.env"); err != nil {
		return nil, "", fmt.Errorf("ошибка загрузки .env: %w", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return nil, "", fmt.Errorf("ошибка подключения к БД: %w", err)
	}

	// проверяем, что коннект рабочий
	if err := db.Ping(); err != nil {
		return nil, "", fmt.Errorf("ping db failed: %w", err)
	}

	return db, connStr, nil
}
