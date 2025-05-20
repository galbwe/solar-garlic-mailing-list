package handlers

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log/slog"
	"net/http"

	"solar-garlic-mailing-list/database"
	"solar-garlic-mailing-list/model"

	"github.com/go-playground/validator/v10"
)

// really a factory that creates a handler function
func CreateEmailHandler(db *sql.DB, validate *validator.Validate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the email from the request body

		var body = &model.CreateEmailBody{}

		err := json.NewDecoder(r.Body).Decode(body)
		if err != nil {
			slog.Info("CreateEmailHandler received invalid request body", "body", *body)
			errorResponse(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		emailAddress := body.Email

		if errs := validate.Var(emailAddress, "required,email"); errs != nil {
			slog.Error("Invalid email address", "email", emailAddress, "errs", errs)
			errorResponse(w, "Invalid body parameter: email. Please use a valid email address.", http.StatusBadRequest)
			return
		}

		email, dbErr := database.CreateEmail(db, emailAddress)

		// will need to consider case where the email already exists in the db

		if dbErr != nil {
			slog.Error("Database error while creating email", "err", dbErr)

			// check if the email already exists
			existing, err := database.ListEmails(db, emailAddress)
			if err != nil {
				slog.Error("Database error while creating email", "err", err)
				errorResponse(w, "Error creating user with that email", http.StatusInternalServerError)
				return
			}

			if len(existing) > 0 {
				slog.Error("Tried to create email address that already exists", "email", email)
				errorResponse(w, "Email address already exists", http.StatusConflict)
				return
			}

			errorResponse(w, "Error creating user with that email", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encodeErr := json.NewEncoder(w).Encode(email)
		if encodeErr != nil {
			slog.Error("Error encoding CreateEmailHandler response body", "err", encodeErr)
			errorResponse(w, "Error encoding response", http.StatusInternalServerError)
			return
		}

	}
}

func ListEmailsHandler(db *sql.DB, validate *validator.Validate) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		email, decodeErr := decodeEmailAddressParam(r)
		if decodeErr != nil {
			slog.Error("Error while decoding query string param", "param", "a", "err", decodeErr)
			errorResponse(w, "Invalid query string param: a. Please use a url-safe base64 encoded email.", http.StatusBadRequest)
			return
		}
		if email != "" {
			// TODO: validate the email
			slog.Info("Querying for emails by address.", "email", email)
			if errs := validate.Var(email, "required,email"); errs != nil {
				slog.Error("Validation error for email", "email", email, "errs", errs)
				errorResponse(w, "Invalid query string param: a. Please use a url-safe base64 encoded email.", http.StatusBadRequest)
				return
			}
		}
		emails, dbErr := database.ListEmails(db, email)
		if dbErr != nil {
			slog.Error("Database error while listing emails", "err", dbErr)
			errorResponse(w, "Error listings emails", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		encodeErr := json.NewEncoder(w).Encode(model.ListEmailsResponse{
			Emails: emails,
		})
		if encodeErr != nil {
			slog.Error("Error encoding ListEmailsHandler response body", "err", encodeErr)
			errorResponse(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

func decodeEmailAddressParam(r *http.Request) (string, error) {
	query := r.URL.Query()

	encodedEmail := query.Get("a")
	email, err := base64.RawURLEncoding.DecodeString(encodedEmail)

	if err != nil {
		return "", err
	}

	return string(email), nil
}

func errorResponse(w http.ResponseWriter, msg string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(model.ErrorResponse{Message: msg})
	if err != nil {
		slog.Error("Could not encode error response", "err", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
