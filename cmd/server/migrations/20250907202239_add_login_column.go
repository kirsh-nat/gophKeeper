package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddLoginColumn, downAddLoginColumn)
}

func upAddLoginColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE users ADD COLUMN login VARCHAR(255);`)
	if err != nil {
		return err
	}
	return nil
}

func downAddLoginColumn(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE users DROP COLUMN login;`)
	if err != nil {
		return err
	}
	return nil
}
