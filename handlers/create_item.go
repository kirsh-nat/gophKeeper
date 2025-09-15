package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"gophkeer/server/internal/app"
	"gophkeer/server/internal/models/item"
	itemservices "gophkeer/server/internal/services/item_services"
	"net/http"
)

func (h *KeeperHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	activeUser, ok := getActiveUser(r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized"))
		return
	}

	var buf bytes.Buffer
	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var userData item.UserDataItem

	if err = json.Unmarshal(buf.Bytes(), &userData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	app.Sugar.Debug("Create user data file: ", userData)

	_, err = itemservices.CreateUserItem(r.Context(), h.db, userData, activeUser.ID)
	if err != nil {
		var undefinedTypeErr *item.UndefinedDataTypeError
		if errors.As(err, &undefinedTypeErr) {
			app.Sugar.Errorw(err.Error(), "event", "create item by type")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		var invalidDataErr *item.DataStrutureError
		if errors.As(err, &invalidDataErr) {
			app.Sugar.Errorw(err.Error(), "event", "invalid data structure to create item")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		app.Sugar.Errorw(err.Error(), "event", "something went wrong with create item")
		h.StatusServerError(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("Item sucessfully created"))
}
