package app

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func setDB() {
	errEnv := godotenv.Load("/opt/gophKeeper/.env")
	if errEnv != nil {
		log.Fatalf("Ошибка загрузки .env: %v", errEnv)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	ConnStr = "host=" + dbHost + " port=" + dbPort + " user=" + dbUser + " password=" + dbPass + " dbname=" + dbName
	var err error
	DB, err = sql.Open("pgx", ConnStr)
	if err != nil {
		panic(err)
	}
}
