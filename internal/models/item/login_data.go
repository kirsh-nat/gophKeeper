package item

import (
	"context"
	"database/sql"
	"encoding/json"
)

// LoginData служит для хранения информации типа логин - пароль
type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Info     string `json:"name"`
}

// CreateLoginItem создает новый элемент типа логин - пароль
//
//   - db - объект базы данных
//   - DataType - тип данных, в данном случае "login_data"
//   - userID - айди владельца элемента
//   - userData - данные для создания элемента
func CreateLoginItem(ctx context.Context, db *sql.DB, DataType string, userID int, userData LoginData) (*Item, error) {
	var item Item

	attrs, err := json.Marshal(userData)
	if err != nil {
		return nil, err
	}

	query := `
        INSERT INTO items (user_id, type, attributes_json)
        VALUES ($1, $2, $3)
        RETURNING id, user_id, type, attributes_json,  created_at, updated_at
    `
	err = db.QueryRowContext(ctx, query, userID, DataType, attrs).
		Scan(&item.ID, &item.UserID, &item.Type, &item.AttributesJSON,
			&item.CreatedAt, &item.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &item, nil
}
