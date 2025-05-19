package handlers

import (
	"database/sql"
	"encoding/json"
	"log/slog"
	"net/http"

	"solar-garlic-mailing-list/database"
	"solar-garlic-mailing-list/model"
)

// really a factory that creates a handler function
func CreateEmailHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the email from the request body

		var body = &model.CreateEmailBody{}

		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			slog.Info("CreateEmailHandler received invalid request body", "body", *body)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		emailAddress := body.Email

		email, dbErr := database.CreateEmail(db, emailAddress)

		// will need to consider case where the email already exists in the db

		if dbErr != nil {
			slog.Error("Database error while creating email", "err", dbErr)

			// check if the email already exists

			http.Error(w, "Error creating user with that email", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		encodeErr := json.NewEncoder(w).Encode(email)
		if encodeErr != nil {
			slog.Error("Error encoding CreateEmailHandler response body", "err", encodeErr)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}

	}
}

func ListEmailsHandler(db *sql.DB) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		emails, dbErr := database.ListEmails(db)
		if dbErr != nil {
			slog.Error("Database error while listing emails", "err", dbErr)
			http.Error(w, "Error listings emails", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")

		encodeErr := json.NewEncoder(w).Encode(model.ListEmailsResponse{
			Emails: emails,
		})
		if encodeErr != nil {
			slog.Error("Error encoding ListEmailsHandler response body", "err", encodeErr)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}
