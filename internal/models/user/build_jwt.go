package user

// func BuildJWTString(UUID string) (string, error) {
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
// 		},
// 		UserID: UUID,
// 	})

// 	tokenString, err := token.SignedString([]byte(secretKey))
// 	if err != nil {
// 		return "", err
// 	}

// 	return tokenString, nil
// }
