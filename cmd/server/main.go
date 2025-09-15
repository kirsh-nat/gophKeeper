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

	db, connStr, err := app.InitDB()
	if err != nil {
		app.Sugar.Fatalw("cannot init db", "err", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			app.Sugar.Errorw("Error closing database", "error", err)
		}
	}()

	if err := migrations.RunMigrations(connStr); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start db")
	}

	handler := handlers.NewHandler(db)

	if err := run(handler); err != nil {
		app.Sugar.Fatalw(err.Error(), "event", "start server")
	}
}

func run(handler *handlers.KeeperHandler) error {

	mux := handlers.Routes(handler)
	fmt.Print("Server started on: ", app.Adress)

	return http.ListenAndServe(app.Adress, mux)
}
