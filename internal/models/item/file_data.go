package item

import (
	"context"
	"database/sql"
)

// CreateFileItem создает новый Item для файла
func CreateFileItem(ctx context.Context, db *sql.DB, userID int, filename string) (*Item, error) {
	var item Item

	query := `
        INSERT INTO items (user_id, type, attributes_json)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, type, attributes_json, created_at, updated_at
    `

	// В attributes_json можно хранить произвольные метаданные (например, оригинальное имя файла)
	attrs := []byte(`{"filename":"` + filename + `"}`)

	err := db.QueryRowContext(ctx, query, userID, TypeFile, attrs).
		Scan(&item.ID, &item.UserID, &item.Type, &item.AttributesJSON,
			&item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
