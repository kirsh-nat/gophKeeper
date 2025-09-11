package attachment

import (
	"context"
	"database/sql"
)

func Insert(ctx context.Context, db *sql.DB, att *Attachment) error {
	query := `
        INSERT INTO attachments (item_id, content_type, size_bytes, storage_key, storage_backend, sha256, chunk_size, chunks_count)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at
    `
	return db.QueryRowContext(
		ctx,
		query,
		att.ItemID,
		att.ContentType,
		att.SizeBytes,
		att.StorageKey,
		att.StorageBackend,
		att.Sha256,
		att.ChunkSize,
		att.ChunksCount,
	).Scan(&att.ID, &att.CreatedAt)
}
