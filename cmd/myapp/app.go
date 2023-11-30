package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
	"net/http"
)

type App struct {
	Router *chi.Mux
	DB     *sqlx.DB
}

func NewApp(db *sqlx.DB) *App {
	app := &App{
		Router: chi.NewRouter(),
		DB:     db,
	}
	app.setupRoutes()
	return app
}

func (app App) setupRoutes() {
	// Configuring routes
	app.Router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome to the App!"))
	})
}

func (app *App) Start(addr string) error {
	return http.ListenAndServe(addr, app.Router)
}
