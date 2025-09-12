package userservices

import (
	"context"
	"database/sql"
	"gophkeer/server/internal/models/user"
)

func FindOne(DB *sql.DB, ctx context.Context, login, password string) (*user.User, error) {
	fuser := &user.User{}
	err := DB.QueryRowContext(ctx, "SELECT id, login, password_salt FROM users WHERE login = $1", login).Scan(&fuser.ID, &fuser.Login, &fuser.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, user.NewAuthorizationError("Cannot find user", err)
		}
		return nil, err
	}
	authorizationError := checkPassword(fuser.Password, password)
	if authorizationError != nil {
		return nil, user.NewAuthorizationError("Cannot find user", authorizationError)
	}

	return fuser, nil
}
