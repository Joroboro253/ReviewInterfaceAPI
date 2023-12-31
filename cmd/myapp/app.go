package main

import (
	"ReviewInterfaceAPI/internal/handlers"
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

func (app *App) Start(addr string) error {
	return http.ListenAndServe(addr, app.Router)
}

func (app *App) setupRoutes() {
	reviewHandler := &handlers.Handler{
		DB: app.DB,
	}
	app.Router.Post("/products/{product_id}/reviews", reviewHandler.CreateReview)
	app.Router.Get("/products/{product_id}/reviews", reviewHandler.GetReviews)
	app.Router.Delete("/products/{product_id}/reviews", reviewHandler.DeleteReviews)
	app.Router.Patch("/products/{product_id}/reviews/{review_id}", reviewHandler.UpdateReviewById)
}
