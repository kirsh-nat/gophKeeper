package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"gophkeer/server/internal/app"
	"gophkeer/server/internal/models/attachment"
	"gophkeer/server/internal/models/item"
	"gophkeer/server/internal/services"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

// UploadFile обрабатывает загрузку файла на сервер.
// Он принимает multipart-форму, содержащую файл, и возвращает
// HTTP-код Accepted (202) сразу после начала обработки.
// После этого он создает Item (статус pending) и возвращает его ID
// в теле ответа.
// Затем он асинхронно обрабатывает файл, сохраняет его в хранилище
// (local storage) и обновляет Attachment в БД.
func (h *KeeperHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	activeUser, ok := h.getUserFromToken(w, r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаем Item (статус pending)
	newItem, err := item.CreateFileItem(r.Context(), h.db, activeUser.ID, header.Filename)
	if err != nil {
		http.Error(w, "failed to create item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Сохраняем во временный файл (только чтобы передать в горутину)
	tmpFile, err := os.CreateTemp("", "upload_*")
	if err != nil {
		http.Error(w, "failed to create temp file", http.StatusInternalServerError)
		return
	}
	defer tmpFile.Close()

	if _, err := io.Copy(tmpFile, file); err != nil {
		http.Error(w, "failed to buffer file", http.StatusInternalServerError)
		return
	}

	// Возвращаем сразу Accepted
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(fmt.Sprintf(`{"item_id":%d}`, newItem.ID)))

	go func(tmpPath string, hFile *multipart.FileHeader, itemID int, db *sql.DB) {
		defer os.Remove(tmpPath)

		storageDir := "storage"
		os.MkdirAll(storageDir, 0755)

		normalizedName := services.NormalizeFilename(hFile.Filename)
		storageKey := fmt.Sprintf("%s/%d_%s", storageDir, itemID, normalizedName)

		src, err := os.Open(tmpPath)
		if err != nil {
			app.Sugar.Errorw("failed to reopen temp file", "err", err)
			return
		}
		defer src.Close()

		out, err := os.Create(storageKey)
		if err != nil {
			app.Sugar.Errorw("failed to create storage file", "err", err)
			return
		}
		defer out.Close()

		hasher := sha256.New()
		size, err := io.Copy(io.MultiWriter(out, hasher), src)
		if err != nil {
			app.Sugar.Errorw("failed to copy to storage", "err", err)
			return
		}

		ct := detectContentType(storageKey, hFile.Header.Get("Content-Type"))

		att := &attachment.Attachment{
			ItemID:         itemID,
			SizeBytes:      size,
			StorageKey:     storageKey,
			StorageBackend: "local",
			Sha256:         hasher.Sum(nil),
			ContentType:    ct,
			Status:         "done",
		}

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := attachment.Update(ctx, db, att); err != nil {
			app.Sugar.Errorw("failed to update attachment", "err", err)
			return
		}

	}(tmpFile.Name(), header, newItem.ID, h.db)
}

// detectContentType возвращает content-type для переданного пути и fallback
func detectContentType(path, fallback string) string {
	f, err := os.Open(path)
	if err != nil {
		return fallback
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, _ := f.Read(buf)
	detected := http.DetectContentType(buf[:n])
	if detected == "application/octet-stream" && fallback != "" {
		return fallback
	}
	return detected
}
