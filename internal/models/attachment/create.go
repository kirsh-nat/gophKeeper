package attachment

import "database/sql"

func CreateAttachment(db *sql.DB, itemID int, contentType string, sizeBytes int64, storageKey string, storageBackend string, sha256 []byte, chunkSize int, chunksCount int) (*Attachment, error) {
	var att Attachment
	query := `
		INSERT INTO attachments (item_id, content_type, size_bytes, storage_key, storage_backend, sha256, chunk_size, chunks_count)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, item_id, content_type, size_bytes, storage_key, storage_backend, sha256, chunk_size, chunks_count, created_at
	`
	err := db.QueryRow(query, itemID, contentType, sizeBytes, storageKey, storageBackend, sha256, chunkSize, chunksCount).
		Scan(&att.ID, &att.ItemID, &att.ContentType, &att.SizeBytes, &att.StorageKey, &att.StorageBackend, &att.Sha256, &att.ChunkSize, &att.ChunksCount, &att.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &att, nil
}
