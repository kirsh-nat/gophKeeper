package user

import (
	"context"
	"database/sql"
)

func GetByID(DB *sql.DB, ctx context.Context, id int) (*User, error) {
	user := &User{}
	err := DB.QueryRowContext(ctx, "SELECT id, login FROM users WHERE id = $1", id).Scan(&user.ID, &user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
