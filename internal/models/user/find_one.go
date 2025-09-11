package user

import (
	"context"
	"database/sql"
)

func FindOne(DB *sql.DB, ctx context.Context, login, password string) (*User, error) {
	user := &User{}
	err := DB.QueryRowContext(ctx, "SELECT id, login, password_salt FROM users WHERE login = $1", login).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewAuthorizationError("Cannot find user", err)
		}
		return nil, err
	}
	authorizationError := checkPassword(user.Password, password)
	if authorizationError != nil {
		return nil, NewAuthorizationError("Cannot find user", authorizationError)
	}

	return user, nil
}
