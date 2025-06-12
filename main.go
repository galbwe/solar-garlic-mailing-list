package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"solar-garlic-mailing-list/database"
	"solar-garlic-mailing-list/email"
	"solar-garlic-mailing-list/handlers"
	"solar-garlic-mailing-list/jobs"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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
	SKIP_MAILING_LIST_VERIFICATION := setConfig("SKIP_MAILING_LIST_VERFICATION", false)
	EMAIL_VERIFCATION_TOKEN_TTL_SECONDS, err := strconv.ParseInt(setConfig("EMAIL_VERIFICATION_TOKEN_TTL_SECONDS", false), 10, 64)
	if err != nil {
		slog.Error("Error reading env var", "err", err, "var", "EMAIL_VERIFICATION_TOKEN_TTL_SECONDS")
		panic("Could not read env var: EMAIL_VERIFICATION_TOKEN_TTL_SECONDS")
	}
	EMAIL_VERIFICATION_SUCCESS_REDIRECT := setConfig("EMAIL_VERIFICATION_SUCCESS_REDIRECT", false)
	EMAIL_VERIFICATION_ENDPOINT := setConfig("EMAIL_VERIFICATION_ENDPOINT", false)
	EMAIL_UNSUBSCRIBE_ENDPOINT := setConfig("EMAIL_UNSUBSCRIBE_ENDPOINT", false)
	ALLOW_ALL_CORS := setConfig("ALLOW_ALL_CORS", false)

	// configure logging
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	slog.Info("Creating sql database")
	db = database.CreateDB(SQLITE_DB)
	database.CreateSchema(db)

	r := chi.NewRouter()

	// request logging middleware
	r.Use(middleware.Logger)

	// CORS middleware
	if ALLOW_ALL_CORS == "true" {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	var validate = validator.New(validator.WithRequiredStructEnabled())

	// v1 routes
	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/emails", handlers.ListEmailsHandler(db, validate))
		r.Post("/emails", handlers.CreateEmailHandler(db, validate, MAIL_CONFIG, SKIP_MAILING_LIST_VERIFICATION, EMAIL_VERIFICATION_ENDPOINT))

		r.Get("/emails/verify", handlers.VerifyEmail(db, validate, MAIL_CONFIG, int(EMAIL_VERIFCATION_TOKEN_TTL_SECONDS), EMAIL_VERIFICATION_SUCCESS_REDIRECT, EMAIL_UNSUBSCRIBE_ENDPOINT))
		r.Get("/emails/unsubscribe", handlers.UnsubscribeEmail(db, validate))

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
