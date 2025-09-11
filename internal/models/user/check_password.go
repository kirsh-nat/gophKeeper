package user

import "golang.org/x/crypto/bcrypt"

func checkPassword(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return NewAuthorizationError("Check password!", err)
	}
	return nil
}
