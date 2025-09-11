package user

import (
	"context"
	"database/sql"
)

// TODO: вынести в отдельный пакет вероятно, чтобы замокать легко было
func GetByPassword(DB *sql.DB, ctx context.Context, password string) (*User, error) {
	user := &User{}
	err := DB.QueryRowContext(ctx, "SELECT id, username, password FROM users WHERE password = $1", password).Scan(&user.ID, &user.Login, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil
		}
		return user, err
	}
	return user, nil
}
