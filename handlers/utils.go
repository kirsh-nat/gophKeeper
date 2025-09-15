package handlers

import (
	"database/sql"
	"fmt"
	"gophkeer/server/internal/app"
	"gophkeer/server/internal/models/user"
	userservices "gophkeer/server/internal/services/user-services"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type dataUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type KeeperHandler struct {
	db *sql.DB
}

const tokenExp = time.Hour * 3

func NewHandler(db *sql.DB) *KeeperHandler {
	return &KeeperHandler{db: db}
}

func (h *KeeperHandler) checkMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method not allowed"))
		return false
	}

	return true
}

func (h *KeeperHandler) StatusBadRequest(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Bad request"))
	return false
}

func (h *KeeperHandler) StatusServerError(w http.ResponseWriter, r *http.Request) bool {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Something went wrong"))
	return false
}

func (h *KeeperHandler) setCookieToken(user *user.User, w http.ResponseWriter) (*user.User, bool) {
	token, err := createToken(user)
	if err != nil {
		return nil, false
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		//Secure:   true, // ⚠️ включать только если у тебя HTTPS
		SameSite: http.SameSiteStrictMode,
	})

	return user, true
}

func createToken(user *user.User) (string, error) {
	token, err := userservices.BuildJWTString(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil

}

func (h *KeeperHandler) getUserFromToken(w http.ResponseWriter, r *http.Request) (*user.User, bool) {
	cookieToken, err := r.Cookie("token")
	if err != nil {
		return nil, false
	}

	userID, err := getUserID(cookieToken.Value)

	if err != nil {
		app.Sugar.Errorw("failed to parse token", "error", err)
		return nil, false
	}

	foundUser, err := userservices.GetByID(h.db, r.Context(), userID)
	if err != nil {
		app.Sugar.Errorw("failed to get user from DB", "error", err, "user_id", userID)
		return nil, false
	}

	return foundUser, true
}

func getUserID(tokenString string) (int, error) {
	claims := &user.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			// проверяем, что используется HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				app.Sugar.Errorf("unexpected signing method: %v", t.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(userservices.SecretKey), nil
		},
	)
	if err != nil {
		app.Sugar.Error("failed to parse token", "error", err)
		return 0, err
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	return claims.UserID, nil
}

// func buildJWTString(UUID int) (string, error) {
// 	claims := Claims{
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
// 		},
// 		UserID: UUID,
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	tokenString, err := token.SignedString([]byte(secretKey))
// 	if err != nil {
// 		return "", err
// 	}
// 	return tokenString, nil
// }
