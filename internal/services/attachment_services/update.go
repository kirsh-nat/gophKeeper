package attachmentservices

import (
	"context"
	"database/sql"
	"fmt"
	"gophkeer/server/internal/models/attachment"
)

// Update обновляет существующий Attachment частично
// Передавай только те поля, которые нужно изменить.
// Например: SizeBytes, StorageKey, Sha256, Status
func Update(ctx context.Context, db *sql.DB, att *attachment.Attachment) error {
	query := `
        UPDATE attachments
        SET
            size_bytes = COALESCE($1, size_bytes),
            storage_key = COALESCE($2, storage_key),
            sha256 = COALESCE($3, sha256),
            status = COALESCE($4, status)
        WHERE id = $5
    `

	_, err := db.ExecContext(
		ctx,
		query,
		nullInt64(att.SizeBytes),
		nullString(att.StorageKey),
		nullBytes(att.Sha256),
		nullString(att.Status),
		att.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update attachment: %w", err)
	}

	return nil
}

// nullInt64 возвращает *int64 для COALESCE (0 → nil)
func nullInt64(v int64) interface{} {
	if v == 0 {
		return nil
	}
	return v
}

// nullString возвращает *string для COALESCE ("" → nil)
func nullString(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

// nullBytes возвращает []byte для COALESCE (nil → nil)
func nullBytes(b []byte) interface{} {
	if b == nil || len(b) == 0 {
		return nil
	}
	return b
}
