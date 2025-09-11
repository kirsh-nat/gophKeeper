package services

import (
	"gophkeer/server/internal/models/user"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// tokenExp задает время истечения срока действия токена (3 дня)
const tokenExp = time.Hour * 72

// secretKey — секретный ключ для подписи JWT
const SecretKey = "supersecretkeyJWT"

func BuildJWTString(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, user.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		UserID: id,
	})

	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
