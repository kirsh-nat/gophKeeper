package handlers

import (
	"encoding/json"
	"gophkeer/server/internal/app"
	userservices "gophkeer/server/internal/services/user-services"
	"net/http"
)

func (h *KeeperHandler) Registration(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	var req dataUser
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.StatusBadRequest(w, r)
		return
	}

	user, err := userservices.CreateUser(h.db, r.Context(), req.Login, req.Password)
	if err != nil {
		app.Sugar.Errorw(err.Error(), "event", "create user")
		h.StatusServerError(w, r)
		return
	}

	user, ok := h.setCookieToken(user, w)
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Something went wrong"))
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("User was successfully created: " + user.Login))
}
