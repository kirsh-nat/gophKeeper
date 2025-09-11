package item

// func TestCreateLoginItem(t *testing.T) {
// 	db, mock, err := sqlmock.New()
// 	require.NoError(t, err)
// 	defer db.Close()

// 	userID := 1
// 	dataType := "login"
// 	loginData := item.LoginData{
// 		Login:    "user123",
// 		Password: "pass123",
// 		Info:     "Test User",
// 	}

// 	// ожидаемый JSON
// 	attrs, _ := json.Marshal(loginData)

// 	// мокаем текущее время
// 	now := time.Now()

// 	rows := sqlmock.NewRows([]string{
// 		"id", "user_id", "type", "attributes_json", "wrapped_data_key", "version", "created_at", "updated_at",
// 	}).AddRow(1, userID, dataType, attrs, []byte{}, 1, now, now)

// 	mock.ExpectQuery(regexp.QuoteMeta(`
//         INSERT INTO items (user_id, type, attributes_json)
//         VALUES ($1, $2, $3)
//         RETURNING id, user_id, type, attributes_json, wrapped_data_key, version, created_at, updated_at
//     `)).
// 		WithArgs(userID, dataType, driver.Value(attrs)).
// 		WillReturnRows(rows)

// 	itemResult, err := item.CreateLoginItem(db, dataType, userID, loginData)
// 	require.NoError(t, err)
// 	require.NotNil(t, itemResult)
// 	require.Equal(t, 1, itemResult.ID)
// 	require.Equal(t, userID, itemResult.UserID)
// 	require.Equal(t, dataType, itemResult.Type)

// 	decoded, err := itemResult.DecodeLoginData()
// 	require.NoError(t, err)
// 	require.Equal(t, loginData.Login, decoded.Login)
// 	require.Equal(t, loginData.Password, decoded.Password)
// 	require.Equal(t, loginData.Info, decoded.Info)
// }
