package handlers

import (
	"net/http"

	"github.com/go-chi/chi"
)

func Routes(handler *KeeperHandler) *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/user/register", http.HandlerFunc(handler.Registration))
	r.Post("/api/user/login", http.HandlerFunc(handler.Authentication))
	r.Post("/createItem", http.HandlerFunc(handler.CreateItem))
	r.Post("/uploadFile", http.HandlerFunc(handler.UploadFile))
	r.Get("/downloadFile", http.HandlerFunc(handler.DownloadFile))

	return r
}
