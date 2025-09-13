package app

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func SetAppConfig() {
	setLogger()
	setDB()
	errEnv := godotenv.Load("/opt/gophKeeper/.env")
	if errEnv != nil {
		log.Fatalf("Ошибка загрузки .env: %v", errEnv)
	}

	Storage = os.Getenv("FILES_PATH")
	Adress = os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
}
