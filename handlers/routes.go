package handlers

import (
	"github.com/go-chi/chi"
)

func Routes(handler *KeeperHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/api/user/register", handler.Registration)
	r.Post("/api/user/login", handler.Authentication)

	r.With(handler.Middleware).Post("/createItem", handler.CreateItem)
	r.With(handler.Middleware).Post("/uploadFile", handler.UploadFile)
	r.With(handler.Middleware).Get("/downloadFile", handler.DownloadFile)

	return r
}
