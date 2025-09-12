package userservices

import (
	"gophkeer/server/internal/models/user"

	"golang.org/x/crypto/bcrypt"
)

func checkPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return user.NewAuthorizationError("Check password!", err)
	}
	return nil
}
