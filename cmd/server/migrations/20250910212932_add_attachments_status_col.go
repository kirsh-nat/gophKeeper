package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddAttachmentsStatusCol, downAddAttachmentsStatusCol)
}

func upAddAttachmentsStatusCol(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE attachments ADD COLUMN status VARCHAR(15);`)
	if err != nil {
		return err
	}
	return nil
}

func downAddAttachmentsStatusCol(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.ExecContext(ctx, `ALTER TABLE attachments DROP COLUMN status;`)
	if err != nil {
		return err
	}
	return nil
}
