package item

import (
	"encoding/json"
	"time"
)

// Типы хранимых данных в таблице
const (
	TypeFile      = "file"
	TypeText      = "text"
	TypeLoginData = "login_data"
)

type UserDataItem struct {
	DataType string          `json:"data_type"`
	Data     json.RawMessage `json:"data"`
}

type Item struct {
	ID             int       `json:"id"`
	UserID         int       `json:"user_id"`
	Type           string    `json:"type"`
	TitleCipher    []byte    `json:"title_cipher"`
	AttributesJSON []byte    `json:"attributes_json"`
	WrappedDataKey []byte    `json:"wrapped_data_key"`
	Version        int       `json:"version"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
