package user

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(DB *sql.DB, ctx context.Context, login, password string) (*User, error) {
	hashPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	result, err := DB.ExecContext(ctx,
		"INSERT INTO users (login, password_salt) VALUES ($1, $2)", login, hashPassword)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return nil, NewUserExistsError("Create user", err)
		}
		return nil, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if rowsAffected == 0 {
		return nil, NewUserExistsError("Create user", err)
	}

	return FindOne(DB, ctx, login, password)
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
