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
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var db *sql.DB

func main() {
	godotenv.Load()

	// configuration variables
	PORT := setConfig("PORT")
	SQLITE_DB := setConfig("SQLITE_DB")

	// configure logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("Creating sql database")
	db = database.CreateDB(SQLITE_DB)
	database.CreateSchema(db)

	r := chi.NewRouter()

	r.Use(middleware.Logger)

	var validate = validator.New(validator.WithRequiredStructEnabled())

	// v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/emails", handlers.ListEmailsHandler(db, validate))
		r.Post("/emails", handlers.CreateEmailHandler(db, validate))

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

	slog.Info("Starting the server", "port", PORT)
	http.ListenAndServe(":"+PORT, r)

}

func setConfig(v string) string {
	if val, ok := os.LookupEnv(v); ok {
		slog.Info("Setting environment variable", v, val)
		return val
	}
	slog.Warn("Environment variable could not be read", "variable", v)
	return ""
}
