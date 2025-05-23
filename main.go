package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"solar-garlic-mailing-list/database"
	"solar-garlic-mailing-list/email"
	"solar-garlic-mailing-list/handlers"
	"solar-garlic-mailing-list/jobs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/robfig/cron/v3"
)

var db *sql.DB

func main() {
	godotenv.Load()

	// configuration variables
	PORT := setConfig("PORT", false)
	SQLITE_DB := setConfig("SQLITE_DB", false)
	DB_BACKUP_SCHEDULE := setConfig("DB_BACKUP_SCHEDULE", false)
	DB_BACKUP_S3_BUCKET := setConfig("DB_BACKUP_S3_BUCKET", false)
	MAIL_CONFIG := email.MailConfig{
		MailFrom: setConfig("MAIL_FROM", false),
		User:     setConfig("AWS_SMTP_USER", false),
		Password: setConfig("AWS_SMTP_PASSWORD", true),
		Host:     setConfig("SMTP_HOST", false),
	}

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
		r.Post("/emails", handlers.CreateEmailHandler(db, validate, MAIL_CONFIG))

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

	// set up cron jobs
	c := cron.New()

	// backup the sqlite database once a day
	if DB_BACKUP_SCHEDULE != "" {
		c.AddFunc(DB_BACKUP_SCHEDULE, func() {
			slog.Info("Backing up database ...")
			timestamp := time.Now().UTC().Unix()
			key := fmt.Sprintf("/mailing-list/backups/db/%v/%v", timestamp, SQLITE_DB)
			jobs.UploadToS3(SQLITE_DB, DB_BACKUP_S3_BUCKET, key)
		})
	}
	c.Start()
	defer c.Stop()

	// run the server
	slog.Info("Starting the server", "port", PORT)
	http.ListenAndServe(":"+PORT, r)

}

func setConfig(v string, secret bool) string {
	if val, ok := os.LookupEnv(v); ok {
		if secret {
			slog.Info("Setting environment variable", v, val[:4]+"***********************")
		} else {
			slog.Info("Setting environment variable", v, val)
		}

		return val
	}
	slog.Warn("Environment variable could not be read", "variable", v)
	return ""
}
