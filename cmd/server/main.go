package main

import (
	"fmt"
	"gophkeer/server/cmd/server/migrations"
	"gophkeer/server/handlers"
	"gophkeer/server/internal/app"
	"net/http"
)

func main() {
	app.SetAppConfig()
	app.Sugar.Info("Start server")
	err := migrations.RunMigrations(app.ConnStr)
	if err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start db")
	}
	handler := handlers.NewHandler(app.DB)

	if err := run(handler); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
	}

	if app.DB != nil {
		defer app.DB.Close()
	}
}

func run(handler *handlers.KeeperHandler) error {

	mux := handlers.Routes(handler)
	fmt.Print("Server started on: ", app.Adress)
	//return http.ListenAndServe("localhost:8080", mux)
	return http.ListenAndServe(app.Adress, mux)
}
