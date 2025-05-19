package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"solar-garlic-mailing-list/database"
	"solar-garlic-mailing-list/handlers"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var db *sql.DB

func main() {
	// configure logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("Creating sql database")
	db = database.CreateDB()
	database.CreateSchema(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	// v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/emails", handlers.ListEmailsHandler(db))
		r.Post("/emails", handlers.CreateEmailHandler(db))

		r.Get("/emails/{id:^[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(fmt.Sprintf("get email with id %q", id)))
		})
		r.Put("/emails/{id:^[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(fmt.Sprintf("update email with id %q", id)))
		})
		r.Delete("/emails/{id:^[0-9]+}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			w.Write([]byte(fmt.Sprintf("delete email with id %q", id)))
		})

	})

	slog.Info("Starting the server on port 8080")
	http.ListenAndServe(":8080", r)

}
