package handlers

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"fmt"
	"gophkeer/server/internal/app"
	"gophkeer/server/internal/models/attachment"
	"gophkeer/server/internal/services"
	attachmentservices "gophkeer/server/internal/services/attachment_services"
	itemservices "gophkeer/server/internal/services/item_services"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"golang.org/x/sync/errgroup"
)

// UploadFile обрабатывает загрузку файла на сервер.
func (h *KeeperHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if !h.checkMethod(w, r, http.MethodPost) {
		return
	}

	activeUser, ok := getActiveUser(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "failed to get file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Создаем Item (pending)
	newItem, err := itemservices.CreateFileItem(r.Context(), h.db, activeUser.ID, header.Filename)
	if err != nil {
		http.Error(w, "failed to create item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Временный файл для асинхронной обработки
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

	w.WriteHeader(http.StatusAccepted)
	_, _ = w.Write([]byte(fmt.Sprintf(`{"item_id":%d}`, newItem.ID)))

	go func(tmpPath string, hFile *multipart.FileHeader, itemID int, db *sql.DB) {
		defer os.Remove(tmpPath)

		g, ctx := errgroup.WithContext(context.Background())

		g.Go(func() error {
			storageDir := "storage"
			_ = os.MkdirAll(storageDir, 0755)

			normalizedName := services.NormalizeFilename(hFile.Filename)
			storageKey := fmt.Sprintf("%s/%d_%s", storageDir, itemID, normalizedName)

			src, err := os.Open(tmpPath)
			if err != nil {
				return fmt.Errorf("failed to reopen temp file: %w", err)
			}
			defer src.Close()

			out, err := os.Create(storageKey)
			if err != nil {
				return fmt.Errorf("failed to create storage file: %w", err)
			}
			defer out.Close()

			hasher := sha256.New()
			size, err := io.Copy(io.MultiWriter(out, hasher), src)
			if err != nil {
				return fmt.Errorf("failed to copy to storage: %w", err)
			}

			ct := detectContentType(storageKey, hFile.Header.Get("Content-Type"))

			att := &attachment.Attachment{
				ItemID:         itemID,
				SizeBytes:      size,
				StorageKey:     storageKey,
				StorageBackend: h.Storage,
				Sha256:         hasher.Sum(nil),
				ContentType:    ct,
				Status:         "done",
			}

			updateCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
			defer cancel()

			if err := attachmentservices.Update(updateCtx, db, att); err != nil {
				return fmt.Errorf("failed to update attachment: %w", err)
			}

			return nil
		})

		if err := g.Wait(); err != nil {
			app.Sugar.Errorw("upload processing failed", "err", err, "itemID", itemID)
		} else {
			app.Sugar.Infow("upload finished successfully", "itemID", itemID)
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
