package user

import (
	"github.com/golang-jwt/jwt/v4"
)

type User struct {
	ID       int    `json:"id"`            // уникальный идентификатор
	Login    string `json:"login"`         // имя пользователя
	Password string `json:"password_salt"` // хэшированный пароль
}

// Claims представляет структуру данных для JWT
type Claims struct {
	jwt.RegisteredClaims     // встроенные стандартные поля JWT
	UserID               int // пользовательский айдишник
}
