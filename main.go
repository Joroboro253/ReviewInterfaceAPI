package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func main() {
	router := chi.NewRouter()

	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr:    ":3000",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println("failed to listen server", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		// Handle Get
		if r.URL.Path == "/foo" {
			// Handle get foo
		}
	}

	if r.Method == http.MethodPost {
		// Handle POST
	}

	w.Write([]byte("Hello, world!"))
}
