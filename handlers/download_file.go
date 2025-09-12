package handlers

import (
	"fmt"
	attachmentservices "gophkeer/server/internal/services/attachment_services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func (h *KeeperHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodGet) {
		return
	}

	fmt.Print("DownloadFile\n")

	_, ok := h.getUserFromToken(w, r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// получаем item_id из query (?id=18)
	itemIDStr := r.URL.Query().Get("id")
	if itemIDStr == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	// достаем метаданные файла из БД
	att, err := attachmentservices.GetByItemID(r.Context(), h.db, itemID)
	if err != nil {
		http.Error(w, "file not found", http.StatusNotFound)
		return
	}

	//TODO: проверяем, что файл принадлежит текущему пользователю
	// it, err := item.GetByID(r.Context(), h.db, itemID)
	// if err != nil || it.UserID != activeUser.ID {
	// 	http.Error(w, "forbidden", http.StatusForbidden)
	// 	return
	// }

	// открываем файл
	f, err := os.Open(att.StorageKey)
	if err != nil {
		http.Error(w, "cannot open file", http.StatusInternalServerError)
		return
	}
	defer f.Close()

	// ставим заголовки
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", filepath.Base(att.StorageKey)))
	w.Header().Set("Content-Type", att.ContentType)
	w.Header().Set("Content-Length", strconv.FormatInt(att.SizeBytes, 10))

	// безопасная отдача с поддержкой Range-запросов
	http.ServeContent(w, r, att.StorageKey, att.CreatedAt, f)
}
