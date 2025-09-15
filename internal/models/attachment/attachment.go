package attachment

import (
	"time"
)

const (
	StatusPending = "pending"
	StatusDone    = "done"
)

type Attachment struct {
	ID             int       `db:"id"`
	ItemID         int       `db:"item_id"`
	ContentType    string    `db:"content_type"`
	SizeBytes      int64     `db:"size_bytes"`
	StorageKey     string    `db:"storage_key"`
	StorageBackend string    `db:"storage_backend"`
	Status         string    `db:"status"`
	Sha256         []byte    `db:"sha256"`
	ChunkSize      int       `db:"chunk_size"`
	ChunksCount    int       `db:"chunks_count"`
	CreatedAt      time.Time `db:"created_at"`
}
