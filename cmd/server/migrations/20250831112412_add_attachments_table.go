package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddAttachmentsTable, downAddAttachmentsTable)
}

func upAddAttachmentsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `
	CREATE TABLE IF NOT EXISTS attachments (
		id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
		item_id INT NOT NULL REFERENCES items(id) ON DELETE CASCADE,
		content_type VARCHAR(255) NOT NULL,
		size_bytes BIGINT NOT NULL,
		storage_key TEXT NOT NULL,
		storage_backend TEXT NOT NULL,
		sha256 BYTEA NOT NULL,
		chunk_size INT NOT NULL,
		chunks_count INT NOT NULL,
		created_at TIMESTAMPTZ NOT NULL DEFAULT now()
	);`)
	if err != nil {
		return err
	}
	return nil
}

func downAddAttachmentsTable(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, "DROP TABLE attachments")
	if err != nil {
		return err
	}
	return nil
}
