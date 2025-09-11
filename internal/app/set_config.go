package app

import (
	"database/sql"
	"flag"
	"os"
)

func SetAppConfig() {
	setLogger()
	// setDBConfig
	flag.StringVar(&ConnStr,
		//TODO: другая база хахаха
		"d", "host=localhost port=5432 user=gophkeeper password=gophkeeper dbname=gophkeeperdb sslmode=disable",
		"Адрес запуска HTTP-сервера",
	)
	if conn := os.Getenv("DATABASE_URI"); conn != "" {
		ConnStr = conn
	}

	var err error
	DB, err = sql.Open("pgx", ConnStr)
	if err != nil {
		panic(err)
	}
}
