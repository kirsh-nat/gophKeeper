package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddItemsTable, downAddItemsTable)
}

func upAddItemsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS items (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		type VARCHAR(255) NOT NULL,
		title_cipher BYTEA,
		attributes_json JSONB DEFAULT NULL,
		version INT NOT NULL DEFAULT 1,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);`)
	if err != nil {
		return err
	}
	return nil
}

func downAddItemsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE items")
	if err != nil {
		return err
	}
	return nil
}
