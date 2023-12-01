package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
)

func (app *App) setupRoutes() {
	// Configuring routes
	app.Router.Get("/", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Welcome to the App!"))
	})
}

func (app *App) loadRoutes() {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//router.Route("/orders", app.loadOrderRoutes)
}

//func (a *App) loadOrderRoutes(router chi.Router) {
//	orderHandler := &handler.Order{
//		Repo: &order.RedisRepo{
//			Client: a.rdb,
//		},
//	}
//
//	router.Post("/", orderHandler.Create)
//	router.Get("/", orderHandler.List)
//	router.Get("/{id}", orderHandler.GetByID)
//	router.Put("/{id}", orderHandler.UpdateByID)
//	router.Delete("/{id}", orderHandler.DeleteByID)
//}
