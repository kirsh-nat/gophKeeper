package attachment

import (
	"context"
	"database/sql"
	"gophkeer/server/internal/app"
)

func GetByItemID(ctx context.Context, db *sql.DB, itemID int) (*Attachment, error) {
	var att Attachment
	query := `
        SELECT id, item_id, content_type, size_bytes, storage_key, storage_backend, sha256, chunk_size, chunks_count, created_at
        FROM attachments
        WHERE item_id = $1 AND status = $2
        LIMIT 1
    `
	err := db.QueryRowContext(ctx, query, itemID, StatusDone).
		Scan(&att.ID, &att.ItemID, &att.ContentType, &att.SizeBytes, &att.StorageKey,
			&att.StorageBackend, &att.Sha256, &att.ChunkSize, &att.ChunksCount, &att.CreatedAt)
	if err != nil {
		app.Sugar.Error(err)
		return nil, err
	}
	return &att, nil
}
