package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gophkeer/server/internal/models/user"
	userservices "gophkeer/server/internal/services/user-services"
	"net/http"
)

func (h *KeeperHandler) Authentication(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		h.StatusBadRequest(w, r)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &dataUser); err != nil {
		h.StatusBadRequest(w, r)
		return
	}

	u, err := userservices.FindOne(h.db, r.Context(), dataUser.Login, dataUser.Password)
	if err != nil {
		var dErr *user.AuthorizationError
		if errors.As(err, &dErr) {
			w.WriteHeader(http.StatusUnauthorized)
			return

		}
		h.StatusServerError(w, r)
		return
	}

	_, ok := h.setCookieToken(u, w)
	if !ok {
		h.StatusServerError(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}
