package app

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

func SetAppConfig() {
	setLogger()
	var envPath string
	flag.StringVar(&envPath,
		"env", "",
		"Адрес env-файла",
	)

	if envAddr, ok := os.LookupEnv("GOPHKEEPER_ENV"); ok {
		envPath = envAddr
	}

	errEnv := godotenv.Load(envPath)
	if errEnv != nil {
		Sugar.Fatal("Fail to load .env: %v", errEnv)
	}

	Storage = os.Getenv("FILES_PATH")
	Adress = os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")
}
