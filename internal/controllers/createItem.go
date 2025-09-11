package controllers

// type KeeperController struct {
// 	ItemsRepo       *repositories.ItemsRepository
// 	AttachmentsRepo *repositories.AttachmentsRepository
// 	Storage         storage.StorageBackend
// }

// type CreateItemInput struct {
// 	UserID         int
// 	Type           string
// 	TitleCipher    []byte
// 	AttributesJSON []byte
// 	WrappedDataKey []byte
// }

// func (c *KeeperController) CreateItem(ctx context.Context, input CreateItemInput) (*models.Item, error) {
// 	return c.ItemsRepo.Create(ctx, input.UserID, input.Type, input.TitleCipher, input.AttributesJSON, input.WrappedDataKey)
// }

// func (c *KeeperController) UploadAttachment(ctx context.Context, itemID int, file multipart.File, header *multipart.FileHeader) (*models.Attachment, error) {
// 	// считаем sha256
// 	h := sha256.New()
// 	size, _ := io.Copy(h, file)
// 	hash := h.Sum(nil)

// 	// заново открываем (нужно, т.к. io.Copy сдвигает указатель)
// 	file.Seek(0, 0)

// 	storageKey := fmt.Sprintf("items/%d/%d_%s", itemID, time.Now().Unix(), header.Filename)

// 	// грузим файл в хранилище
// 	if err := c.Storage.Upload(ctx, "files", storageKey, file, header.Size, header.Header.Get("Content-Type")); err != nil {
// 		return nil, err
// 	}

// 	// сохраняем метаданные в БД
// 	return c.AttachmentsRepo.Create(ctx, itemID, header.Header.Get("Content-Type"), size, storageKey, "minio", hash, int(size), 1)
// }
