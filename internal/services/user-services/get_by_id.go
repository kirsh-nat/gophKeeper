package userservices

import (
	"context"
	"database/sql"
	"gophkeer/server/internal/models/user"
)

func GetByID(DB *sql.DB, ctx context.Context, id int) (*user.User, error) {
	user := &user.User{}
	err := DB.QueryRowContext(ctx, "SELECT id, login FROM users WHERE id = $1", id).Scan(&user.ID, &user.Login)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
