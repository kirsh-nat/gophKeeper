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
# Показать список команд 
./client-linux -help



