package attachmentservices

import (
	"context"
	"database/sql"
	"gophkeer/server/internal/app"
	"gophkeer/server/internal/models/attachment"
)

func GetByItemID(ctx context.Context, db *sql.DB, itemID int) (*attachment.Attachment, error) {
	var att attachment.Attachment
	query := `
        SELECT id, item_id, content_type, size_bytes, storage_key, storage_backend, sha256, chunk_size, chunks_count, created_at
        FROM attachments
        WHERE item_id = $1 AND status = $2
        LIMIT 1
    `
	err := db.QueryRowContext(ctx, query, itemID, attachment.StatusDone).
		Scan(&att.ID, &att.ItemID, &att.ContentType, &att.SizeBytes, &att.StorageKey,
			&att.StorageBackend, &att.Sha256, &att.ChunkSize, &att.ChunksCount, &att.CreatedAt)
	if err != nil {
		app.Sugar.Error(err)
		return nil, err
	}
	return &att, nil
}
