package app

import (
	"flag"
	"os"

	"github.com/joho/godotenv"
)

// SetAppConfig загружает конфигурационный файл, настройки логгера, и возвращает
// путь для хранения файлов и адрес запуска веб-приложения.
//
// Функция возвращает ошибку, если файл конфигурации не существует.
// В противном случае функция возвращает путь к файлу конфигурации.
//
// Env-файл по умолчанию находится в директории запуска приложения.
// Env-файл может быть переопределен с помощью переменной окружения GOPHKEEPER_ENV.
func SetAppConfig() (string, string) {
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

	storage := os.Getenv("FILES_PATH")
	adress := os.Getenv("APP_HOST") + ":" + os.Getenv("APP_PORT")

	return storage, adress

}
