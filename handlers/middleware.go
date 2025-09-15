package handlers

import (
	"context"
	"gophkeer/server/internal/models/user"
	"net/http"
)

// Middleware для проверки токена и установки пользователя в контекст
func (h *KeeperHandler) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := h.getUserFromToken(w, r)
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// кладем пользователя в контекст
		ctx := context.WithValue(r.Context(), "activeUser", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Хелпер для извлечения пользователя из контекста
func getActiveUser(r *http.Request) (*user.User, bool) {
	user, ok := r.Context().Value("activeUser").(*user.User)
	return user, ok
}
