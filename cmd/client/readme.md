** СБОРКА ПРИЛОЖЕНИЯ ДЛЯ РАЗНЫХ ПЛАТФОРМ *
# Linux
GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=1.0.0 -X 'main.BuildDate=$(date)'" -o client-linux
# macOS
GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.Version=1.0.0 -X 'main.BuildDate=$(date)'" -o client-mac
# Windows
GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=1.0.0 -X 'main.BuildDate=$(date)'" -o client.exe

** ПРИМЕРЫ ЗАПУСКА КЛИЕНТА*
# Показать версию
./client-linux -version
# Зарегистрироваться
./client-linux -action register -user alice -pass 1234
# Логин и создать item
./client-linux -action create -user alice -pass 1234 -item "My File"
# Загрузить файл
./client-linux -action upload -user alice -pass 1234 -file "/home/alice/file.pdf"
# Скачать файл
./client-linux -action download -user alice -pass 1234 -id 1 -file "downloaded.pdf"


