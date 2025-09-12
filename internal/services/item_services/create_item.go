package itemservices

import (
	"context"
	"database/sql"
	"encoding/json"
	"gophkeer/server/internal/models/item"
)

// CreateUserItem создает новый элемент для пользователя
//
//   - ctx - контекст
//   - db - объект базы данных
//   - userData - данные для создания элемента
//   - userID - айди владельца элемента
//
// Возвращает новый элемент, если он успешно создан, или ошибку.
func CreateUserItem(ctx context.Context, db *sql.DB, userData item.UserDataItem, userID int) (*item.Item, error) {

	newItem := &item.Item{}
	switch userData.DataType {
	case item.TypeLoginData:
		var userLoginData item.LoginData
		if err := json.Unmarshal(userData.Data, &userLoginData); err != nil {
			return nil, item.NewDataStrutureError("Can't read data for create login-password data item", err)
		}

		newItem, err := CreateLoginItem(ctx, db, userData.DataType, userID, userLoginData)
		if err != nil {
			return nil, item.NewDataStrutureError("Can't create login-password data item", err)
		}
		return newItem, nil

	}

	return newItem, nil
}
