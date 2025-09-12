package itemservices

import (
	"context"
	"database/sql"
	"gophkeer/server/internal/models/item"
)

// CreateFileItem создает новый Item для файла
func CreateFileItem(ctx context.Context, db *sql.DB, userID int, filename string) (*item.Item, error) {
	var newitem item.Item

	query := `
        INSERT INTO items (user_id, type, attributes_json)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, type, attributes_json, created_at, updated_at
    `

	// В attributes_json можно хранить произвольные метаданные (например, оригинальное имя файла)
	attrs := []byte(`{"filename":"` + filename + `"}`)

	err := db.QueryRowContext(ctx, query, userID, item.TypeFile, attrs).
		Scan(&newitem.ID, &newitem.UserID, &newitem.Type, &newitem.AttributesJSON,
			&newitem.CreatedAt, &newitem.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &newitem, nil
}
