package app

import (
	"flag"
	"os"
)

func setDB() {
	//TODO: вынести отдельно !!!!
	flag.StringVar(&ConnStr,
		"d", "host=localhost port=5432 user=gophkeeper password=gophkeeper dbname=gophkeeperdb sslmode=disable",
		"Адрес запуска HTTP-сервера",
	)
	flag.Parse()

	if conn := os.Getenv("DATABASE_URI"); conn != "" {
		ConnStr = conn
	}
}
